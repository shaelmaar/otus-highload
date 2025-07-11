// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package pg

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	UserCreate(ctx context.Context, arg UserCreateParams) error
	UserGetByID(ctx context.Context, id uuid.UUID) (User, error)
	UserTokenCreate(ctx context.Context, arg UserTokenCreateParams) (int64, error)
	UserTokenDeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

var _ Querier = (*Queries)(nil)
