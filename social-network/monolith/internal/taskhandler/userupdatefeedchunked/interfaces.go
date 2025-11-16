package userupdatefeedchunked

import (
	"context"

	"github.com/google/uuid"
)

type Updater interface {
	UpdateUserFeedChunked(ctx context.Context, userIDs []uuid.UUID) error
}
