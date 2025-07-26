package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByFirstNameLastName(ctx context.Context, firstName string, lastName string) ([]User, error)
	Create(ctx context.Context, user User) error
	DeleteUserTokens(ctx context.Context, userID uuid.UUID) error
	CreateUserToken(ctx context.Context, token UserToken) (int64, error)
	MassCreate(ctx context.Context, users []User) error
	WithTx(tx transaction.Tx) UserRepository
}

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type User struct {
	ID           uuid.UUID
	PasswordHash string
	FirstName    string
	SecondName   string
	BirthDate    time.Time
	Gender       Gender
	Biography    string
	City         string
}

type UserToken struct {
	ID        int64
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
