package friend

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	repo domain.FriendRepository
}

func New(repo domain.FriendRepository) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repo is nil")
	}

	return &UseCases{
		repo: repo,
	}, nil
}
