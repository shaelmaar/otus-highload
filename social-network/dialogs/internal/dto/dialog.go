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
