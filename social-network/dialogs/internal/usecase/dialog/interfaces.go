package dialog

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
)

type KafkaProducer interface {
	MessageCreated(ctx context.Context, event dto.MessageCreatedEvent) error
	MessagesRead(ctx context.Context, event dto.MessagesReadEvent) error
}
