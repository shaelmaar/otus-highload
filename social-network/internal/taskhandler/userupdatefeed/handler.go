package userupdatefeed

import (
	"context"
	"errors"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Handler struct {
	updater Updater
}

func New(updater Updater) (*Handler, error) {
	if utils.IsNil(updater) {
		return nil, errors.New("updater is nil")
	}

	return &Handler{
		updater: updater,
	}, nil
}

func (u *Handler) Handle(ctx context.Context, task dto.UserFeedUpdateTask) error {
	err := u.updater.UpdateUserFeed(ctx, task.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user feed: %w", err)
	}

	return nil
}
