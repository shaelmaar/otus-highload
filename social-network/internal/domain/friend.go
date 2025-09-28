package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FriendRepository interface {
	Create(ctx context.Context, friend Friend) error
	Delete(ctx context.Context, friend Friend) error
}

type Friend struct {
	UserID    uuid.UUID
	FriendID  uuid.UUID
	CreatedAt time.Time
}
