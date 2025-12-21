package dto

import "github.com/google/uuid"

type UnreadDialogMessagesIncrement struct {
	DialogID       string
	MessageID      uint64
	RecipientID    uuid.UUID
	IdempotencyKey string
}

type UnreadDialogMessagesDecrement struct {
	DialogID       string
	MessageIDs     []uint64
	RecipientID    uuid.UUID
	IdempotencyKey string
}

type UnreadDialogMessagesIncremented struct {
	DialogID  string
	MessageID uint64
	Success   bool
}

type UnreadDialogMessagesDecremented struct {
	DialogID   string
	MessageIDs []uint64
	Success    bool
}
