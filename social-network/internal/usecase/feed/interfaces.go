package feed

import (
	"context"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

type Service interface {
	UpdateUserFeedCache(ctx context.Context, userID uuid.UUID) error
}

type UserFeedUpdateProducer interface {
	Publish(ctx context.Context, task dto.UserFeedUpdateTask) error
}
