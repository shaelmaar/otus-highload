package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Repository struct {
	db         pg.QuerierTX
	_replicaDB pg.QuerierTX
}

func New(
	db, replicaDB pg.QuerierTX,
) (*Repository, error) {
	if utils.IsNil(db) {
		return nil, errors.New("db is nil")
	}

	if utils.IsNil(replicaDB) {
		return nil, errors.New("replica db is nil")
	}

	return &Repository{
		db:         db,
		_replicaDB: replicaDB,
	}, nil
}

func (r *Repository) Create(ctx context.Context, post domain.Post) error {
	err := r.db.PostCreate(ctx, pg.PostCreateParams{
		ID:           post.ID,
		Content:      post.Content,
		AuthorUserID: post.AuthorUserID,
		CreatedAt:    post.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create post in db: %w", err)
	}

	return nil
}

func (r *Repository) PostLockByID(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	var out domain.Post

	row, err := r.db.PostGetWithLockByID(ctx, id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return out, domain.ErrPostNotFound
	case err != nil:
		return out, fmt.Errorf("failed to get post by id from db: %w", err)
	}

	return parsePostRow(row), nil
}

func (r *Repository) Update(ctx context.Context, post domain.Post) error {
	err := r.db.PostUpdate(ctx, pg.PostUpdateParams{
		Content: post.Content,
		UpdatedAt: pgtype.Timestamptz{
			Time:             post.UpdatedAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
		ID: post.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to update post in db: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.PostDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete post in db: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	var out domain.Post

	row, err := r.db.PostGetByID(ctx, id)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return out, domain.ErrPostNotFound
	case err != nil:
		return out, fmt.Errorf("failed to get post by id from db: %w", err)
	}

	return parsePostRow(row), nil
}

func (r *Repository) GetLastPostsByUserIDs(
	ctx context.Context, input dto.GetLastPostsByUserIDs) ([]domain.Post, error) {
	rows, err := r.db.LastPostsByUserIDsWithOffsetLimit(ctx, pg.LastPostsByUserIDsWithOffsetLimitParams{
		UserIds: utils.MapSlice(input.UserIDs, func(id uuid.UUID) string {
			return id.String()
		}),
		Offset: int32(input.Offset), //nolint:gosec // здесь не будет значения > 1000.
		Limit:  int32(input.Limit),  //nolint:gosec // здесь не будет значения > 1000.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get last posts by user ids from db: %w", err)
	}

	return utils.MapSlice(rows, parsePostRow), nil
}

func (r *Repository) WithTx(tx transaction.Tx) domain.PostRepository {
	return &Repository{
		db:         r.db.WithTx(tx),
		_replicaDB: nil,
	}
}

func (r *Repository) Slave() domain.PostSlaveRepository {
	return &Repository{
		db:         r._replicaDB,
		_replicaDB: nil,
	}
}

func parsePostRow(row pg.Post) domain.Post {
	return domain.Post{
		ID:           row.ID,
		Content:      row.Content,
		AuthorUserID: row.AuthorUserID,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
