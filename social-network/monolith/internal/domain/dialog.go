package domain

import (
	"time"
)

type DialogMessage struct {
	From      string
	To        string
	Text      string
	CreatedAt time.Time
}
