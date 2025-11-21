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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var out domain.User

	row, err := r.db.UserGetByID(ctx, id)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return out, domain.ErrUserNotFound
	case err != nil:
		return out, fmt.Errorf("failed to get user by id from db: %w", err)
	}

	return parseUserRow(row), nil
}

func (r *Repository) GetByFirstNameLastName(ctx context.Context, firstName, lastName string) ([]domain.User, error) {
	rows, err := r.db.UsersGetByFirstNameSecondName(ctx, pg.UsersGetByFirstNameSecondNameParams{
		FirstName:  firstName,
		SecondName: lastName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get users by first and last name from db: %w", err)
	}

	return utils.MapSlice(rows, parseUserRow), nil
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

func (r *Repository) MassCreate(ctx context.Context, users []domain.User) error {
	_, err := r.db.UsersMassCreate(ctx, utils.MapSlice(users, func(u domain.User) pg.UsersMassCreateParams {
		return pg.UsersMassCreateParams{
			ID:           u.ID,
			PasswordHash: u.PasswordHash,
			FirstName:    u.FirstName,
			SecondName:   u.SecondName,
			BirthDate: pgtype.Date{
				Time:             u.BirthDate,
				InfinityModifier: pgtype.Finite,
				Valid:            !u.BirthDate.IsZero(),
			},
			Gender:    pg.Gender(u.Gender),
			Biography: u.Biography,
			City:      u.City,
		}
	}))
	if err != nil {
		return fmt.Errorf("failed to mass create users in db: %w", err)
	}

	return nil
}

func (r *Repository) WithTx(tx transaction.Tx) domain.UserRepository {
	return &Repository{
		db:         r.db.WithTx(tx),
		_replicaDB: nil,
	}
}

func (r *Repository) Slave() domain.UserSlaveRepository {
	return &Repository{
		db:         r._replicaDB,
		_replicaDB: nil,
	}
}

func parseUserRow(row pg.User) domain.User {
	return domain.User{
		ID:           row.ID,
		PasswordHash: row.PasswordHash,
		FirstName:    row.FirstName,
		SecondName:   row.SecondName,
		BirthDate:    row.BirthDate.Time,
		Gender:       domain.Gender(row.Gender),
		Biography:    row.Biography,
		City:         row.City,
	}
}
