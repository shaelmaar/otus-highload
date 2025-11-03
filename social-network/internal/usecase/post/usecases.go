package post

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

const (
	userFeedChunkSize = 50
)

type UseCases struct {
	repo                    domain.PostRepository
	friendRepo              domain.FriendRepository
	feedService             FeedService
	userFeedChunkedProducer UserFeedChunkedProducer
	postCreatedProducer     PostCreatedChunkedProducer
	tx                      TxExecutor
}

func New(
	repo domain.PostRepository, friendRepo domain.FriendRepository,
	feedService FeedService,
	userFeedChunkedProducer UserFeedChunkedProducer,
	postCreatedProducer PostCreatedChunkedProducer,
	tx TxExecutor,
) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repo is nil")
	}

	if utils.IsNil(friendRepo) {
		return nil, errors.New("friend repo is nil")
	}

	if utils.IsNil(feedService) {
		return nil, errors.New("feed service is nil")
	}

	if utils.IsNil(userFeedChunkedProducer) {
		return nil, errors.New("user feed chunked producer is nil")
	}

	if utils.IsNil(postCreatedProducer) {
		return nil, errors.New("post created chunked producer is nil")
	}

	if utils.IsNil(tx) {
		return nil, errors.New("tx is nil")
	}

	return &UseCases{
		repo:                    repo,
		friendRepo:              friendRepo,
		feedService:             feedService,
		userFeedChunkedProducer: userFeedChunkedProducer,
		postCreatedProducer:     postCreatedProducer,
		tx:                      tx,
	}, nil
}
