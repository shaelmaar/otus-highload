package friend

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	repo                 domain.FriendRepository
	userFeedTaskProducer UserFeedTaskProducer
}

func New(
	repo domain.FriendRepository,
	userFeedTaskProducer UserFeedTaskProducer,
) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repo is nil")
	}

	if utils.IsNil(userFeedTaskProducer) {
		return nil, errors.New("user feed task producer is nil")
	}

	return &UseCases{
		repo:                 repo,
		userFeedTaskProducer: userFeedTaskProducer,
	}, nil
}
