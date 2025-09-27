package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

type PostRepository interface {
	Create(ctx context.Context, post Post) error
	Update(ctx context.Context, post Post) error
	Delete(ctx context.Context, id uuid.UUID) error
	PostLockByID(ctx context.Context, id uuid.UUID) (Post, error)
	WithTx(tx transaction.Tx) PostRepository
	Slave() PostSlaveRepository
}

type PostSlaveRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Post, error)
}

type Post struct {
	ID           uuid.UUID
	Content      string
	AuthorUserID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
