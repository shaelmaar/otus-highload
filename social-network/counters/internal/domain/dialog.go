package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type DialogRepository interface {
	CountUnreadMessages(ctx context.Context, key UnreadDialogMessageCountKey) (int64, error)
	IncrementUnreadMessages(ctx context.Context, key UnreadDialogMessageCountKey, idempotencyKey string) error
	DecrementUnreadMessages(
		ctx context.Context, key UnreadDialogMessageCountKey, decrBy int, idempotencyKey string) error
	Slave() DialogSlaveRepository
}

type DialogSlaveRepository interface {
	CountUnreadMessages(ctx context.Context, key UnreadDialogMessageCountKey) (int64, error)
}

const (
	unreadDialogMessageCountTemplate = `unread:%s_%s`
)

type UnreadDialogMessageCountKey struct {
	DialogID    string
	RecipientID uuid.UUID
}

func (k UnreadDialogMessageCountKey) String() string {
	return fmt.Sprintf(unreadDialogMessageCountTemplate, k.DialogID, k.RecipientID.String())
}
