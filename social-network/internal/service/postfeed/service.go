package postfeed

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

const (
	defaultCachedPostsAmount = 100
)

type Service struct {
	repo       domain.PostRepository
	friendRepo domain.FriendRepository
	cache      Cache
	logger     *zap.Logger
	// какое кол-во последних постов можно получить из кэша.
	cachedPostsAmount int
}

func NewService(
	repo domain.PostRepository,
	friendRepo domain.FriendRepository,
	cache Cache,
	logger *zap.Logger,
) (*Service, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repo is nil")
	}

	if utils.IsNil(friendRepo) {
		return nil, errors.New("friend repo is nil")
	}

	if utils.IsNil(cache) {
		return nil, errors.New("cache is nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	return &Service{
		repo:              repo,
		friendRepo:        friendRepo,
		cache:             cache,
		logger:            logger,
		cachedPostsAmount: defaultCachedPostsAmount,
	}, nil
}

func (s *Service) GetUserFeed(ctx context.Context, input dto.GetPostFeed) ([]domain.Post, error) {
	// если в кэше данных по запросу быть не может - выгружаем из бд.
	if input.Offset+input.Limit > s.cachedPostsAmount {
		return s.getUserFeedFromDB(ctx, input.UserID, input.Offset, input.Limit)
	}

	posts, got, err := s.cache.GetUserFeed(ctx, input)
	if err != nil {
		s.logger.Error("failed to get user feed from cache", zap.Error(err))

		return s.getUserFeedFromDB(ctx, input.UserID, input.Offset, input.Limit)
	}

	if got {
		return posts, nil
	}

	posts, err = s.getUserFeedFromDB(ctx, input.UserID, 0, s.cachedPostsAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to get user feed from db: %w", err)
	}

	err = s.cache.SetUserFeed(ctx, input.UserID, posts)
	if err != nil {
		return nil, fmt.Errorf("failed to set posts in cache: %w", err)
	}

	start := input.Offset
	end := start + input.Limit

	return utils.SafeSliceRange(posts, start, end), nil
}

func (s *Service) UpdateUserFeedCache(ctx context.Context, userID uuid.UUID) error {
	exists, err := s.cache.UserFeedExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user feed existence: %w", err)
	}

	if !exists {
		return nil
	}

	posts, err := s.getUserFeedFromDB(ctx, userID, 0, s.cachedPostsAmount)
	if err != nil {
		return fmt.Errorf("failed to get user feed from db: %w", err)
	}

	err = s.cache.SetUserFeed(ctx, userID, posts)
	if err != nil {
		return fmt.Errorf("failed to set posts in cache: %w", err)
	}

	return nil
}

func (s *Service) getUserFeedFromDB(ctx context.Context, userID uuid.UUID, offset, limit int) ([]domain.Post, error) {
	friendIDs, err := s.friendRepo.Slave().GetUserFriendIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friend ids: %w", err)
	}

	posts, err := s.repo.Slave().GetLastPostsByUserIDs(ctx, dto.GetLastPostsByUserIDs{
		UserIDs: friendIDs,
		Offset:  offset,
		Limit:   limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get last posts: %w", err)
	}

	return posts, nil
}
