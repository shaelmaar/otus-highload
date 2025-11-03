package post

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

type TxExecutor interface {
	Exec(
		ctx context.Context,
		f func(ctx context.Context, tx transaction.Tx) error,
		rollbackFn func(ctx context.Context),
	) error
}

type FeedService interface {
	GetUserFeed(ctx context.Context, input dto.GetPostFeed) ([]domain.Post, error)
}

type UserFeedChunkedProducer interface {
	Publish(ctx context.Context, task dto.UserFeedChunkedUpdateTask) error
}

type PostCreatedChunkedProducer interface {
	Publish(ctx context.Context, task dto.PostCreatedChunkedTask) error
}
