package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DialogRepository interface {
	CreateDialogMessage(ctx context.Context, message DialogMessage) (uint64, error)
	GetMessagesByDialog(
		ctx context.Context, dialogID string) ([]DialogMessage, error)
	MarkMessagesAsReading(
		ctx context.Context, dialogID string, readerID uuid.UUID, messageID uint64) ([]uint64, error)
	UpdateMessageStateFrom(
		ctx context.Context, dialogID string, messageID uint64, fromState, toState DialogMessageState,
	) (bool, error)
	UpdateMessagesStateFrom(
		_ context.Context, dialogID string, messageIDs []uint64, fromState, toState DialogMessageState,
	) ([]uint64, error)
}

type DialogMessageState string

const (
	DialogMessageStateSending DialogMessageState = "sending"
	DialogMessageStateFailed  DialogMessageState = "failed"
	DialogMessageStateSent    DialogMessageState = "sent"
	DialogMessageStateReading DialogMessageState = "reading"
	DialogMessageStateRead    DialogMessageState = "read"
)

type DialogMessage struct {
	ID        uint64
	From      uuid.UUID
	To        uuid.UUID
	DialogID  string
	Text      string
	State     DialogMessageState
	CreatedAt time.Time
}
