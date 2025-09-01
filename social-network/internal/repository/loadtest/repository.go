package loadtest

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Repository struct {
	db pg.QuerierTX
}

func New(
	db pg.QuerierTX,
) (*Repository, error) {
	if utils.IsNil(db) {
		return nil, errors.New("db is nil")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Insert(ctx context.Context, id uuid.UUID, value string) error {
	err := r.db.LoadTestInsert(ctx, pg.LoadTestInsertParams{
		ID:    id,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to insert load test value in db: %w", err)
	}

	return nil
}
