package dialog

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/dto"
)

type KafkaProducer interface {
	UnreadMessagesIncremented(ctx context.Context, event dto.UnreadDialogMessagesIncremented) error
	UnreadMessagesDecremented(ctx context.Context, event dto.UnreadDialogMessagesDecremented) error
}
