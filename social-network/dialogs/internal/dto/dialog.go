package dto

import (
	"time"

	"github.com/google/uuid"
)

type DialogCreateMessage struct {
	From uuid.UUID
	To   uuid.UUID
	Text string
	Time time.Time
}

type DialogMessagesListGet struct {
	From uuid.UUID
	To   uuid.UUID
}

type ReadMessages struct {
	From      uuid.UUID
	To        uuid.UUID
	MessageID uint64
}

type MessageCreatedEvent struct {
	MessageID uint64
	DialogID  string
	From      uuid.UUID
	To        uuid.UUID
}

type MessagesReadEvent struct {
	MessageIDs []uint64
	DialogID   string
	To         uuid.UUID
}
