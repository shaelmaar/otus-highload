package friend

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

type UserFeedTaskProducer interface {
	Publish(ctx context.Context, task dto.UserFeedUpdateTask) error
}
