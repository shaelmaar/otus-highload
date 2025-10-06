package userupdatefeed

import (
	"context"

	"github.com/google/uuid"
)

type Updater interface {
	UpdateUserFeed(ctx context.Context, userID uuid.UUID) error
}
