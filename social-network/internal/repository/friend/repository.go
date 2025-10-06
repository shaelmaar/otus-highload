package friend

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
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

func (r *Repository) Create(ctx context.Context, friend domain.Friend) error {
	err := r.db.FriendCrete(ctx, pg.FriendCreteParams{
		UserID:    friend.UserID,
		FriendID:  friend.FriendID,
		CreatedAt: friend.CreatedAt,
	})

	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr) && pgErr.ConstraintName == "friend_user_id_fkey":
		return domain.ErrUserNotFound
	case errors.As(err, &pgErr) && pgErr.ConstraintName == "friend_friend_id_fkey":
		return domain.ErrFriendNotFound
	case err != nil:
		return fmt.Errorf("failed to create frined in db: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, friend domain.Friend) error {
	err := r.db.FriendDelete(ctx, pg.FriendDeleteParams{
		UserID:   friend.UserID,
		FriendID: friend.FriendID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete friend in db: %w", err)
	}

	return nil
}

func (r *Repository) GetUserFriendIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	friendIDs, err := r.db.FriendIDsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user friend ids from db: %w", err)
	}

	return friendIDs, nil
}

func (r *Repository) GetFriendUserIDs(ctx context.Context, friendID uuid.UUID) ([]uuid.UUID, error) {
	userIDs, err := r.db.UserIDsByFriendID(ctx, friendID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids by friend from db: %w", err)
	}

	return userIDs, nil
}

func (r *Repository) Slave() domain.FriendSlaveRepository {
	return &Repository{
		db:         r._replicaDB,
		_replicaDB: nil,
	}
}
