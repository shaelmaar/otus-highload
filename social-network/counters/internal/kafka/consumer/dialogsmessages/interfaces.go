package dialogsmessages

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/dto"
)

type DialogMessagesCounterUseCases interface {
	IncrementUnreadMessages(ctx context.Context, input dto.UnreadDialogMessagesIncrement) error
	DecrementUnreadMessages(ctx context.Context, input dto.UnreadDialogMessagesDecrement) error
}
