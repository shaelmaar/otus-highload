package userupdatefeedchunked

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

func (u *Handler) Handle(ctx context.Context, task dto.UserFeedChunkedUpdateTask) error {
	err := u.updater.UpdateUserFeedChunked(ctx, task.UserIDs)
	if err != nil {
		return fmt.Errorf("failed to update user feed chunked: %w", err)
	}

	return nil
}
