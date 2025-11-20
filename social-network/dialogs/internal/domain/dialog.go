package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DialogRepository interface {
	CreateDialogMessage(ctx context.Context, message DialogMessage) error
	GetMessagesByDialog(
		ctx context.Context, dialogID string) ([]DialogMessage, error)
}

type DialogMessage struct {
	From      uuid.UUID
	To        uuid.UUID
	DialogID  string
	Text      string
	CreatedAt time.Time
}
