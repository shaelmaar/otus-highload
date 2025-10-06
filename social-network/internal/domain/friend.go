package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FriendRepository interface {
	Create(ctx context.Context, friend Friend) error
	Delete(ctx context.Context, friend Friend) error
	Slave() FriendSlaveRepository
}

type FriendSlaveRepository interface {
	GetUserFriendIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetFriendUserIDs(ctx context.Context, friendID uuid.UUID) ([]uuid.UUID, error)
}

type Friend struct {
	UserID    uuid.UUID
	FriendID  uuid.UUID
	CreatedAt time.Time
}
