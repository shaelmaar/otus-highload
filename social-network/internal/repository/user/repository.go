package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Repository struct {
	db pg.QuerierTX
}

func New(db pg.QuerierTX) (*Repository, error) {
	if utils.IsNil(db) {
		return nil, errors.New("db is nil")
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var out domain.User

	row, err := r.db.UserGetByID(ctx, id)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return out, domain.ErrUserNotFound
	case err != nil:
		return out, fmt.Errorf("failed to get user by id from db: %w", err)
	}

	return domain.User{
		ID:           row.ID,
		PasswordHash: row.PasswordHash,
		FirstName:    row.FirstName,
		SecondName:   row.SecondName,
		BirthDate:    row.BirthDate.Time,
		Gender:       domain.Gender(row.Gender),
		Biography:    row.Biography,
		City:         row.City,
	}, nil
}

func (r *Repository) Create(ctx context.Context, user domain.User) error {
	err := r.db.UserCreate(ctx, pg.UserCreateParams{
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		SecondName:   user.SecondName,
		BirthDate: pgtype.Date{
			Time:             user.BirthDate,
			InfinityModifier: pgtype.Finite,
			Valid:            !user.BirthDate.IsZero(),
		},
		Gender:    pg.Gender(user.Gender),
		Biography: user.Biography,
		City:      user.City,
	})
	if err != nil {
		return fmt.Errorf("failed to create user in db: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUserTokens(ctx context.Context, userID uuid.UUID) error {
	err := r.db.UserTokenDeleteByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user token by user id in db: %w", err)
	}

	return nil
}

func (r *Repository) CreateUserToken(ctx context.Context, token domain.UserToken) (int64, error) {
	id, err := r.db.UserTokenCreate(ctx, pg.UserTokenCreateParams{
		UserID: token.UserID,
		Token:  token.Token,
		ExpiresAt: pgtype.Timestamptz{
			Time:             token.ExpiresAt,
			InfinityModifier: pgtype.Finite,
			Valid:            !token.ExpiresAt.IsZero(),
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create user token in db: %w", err)
	}

	return id, nil
}

func (r *Repository) WithTx(tx transaction.Tx) domain.UserRepository {
	return &Repository{
		db: r.db.WithTx(tx),
	}
}
