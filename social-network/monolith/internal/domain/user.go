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
	MassCreate(ctx context.Context, users []User) error
	WithTx(tx transaction.Tx) UserRepository
	Slave() UserSlaveRepository
}

type UserSlaveRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByFirstNameLastName(ctx context.Context, firstName string, lastName string) ([]User, error)
}

type LoadTestRepository interface {
	Insert(ctx context.Context, id uuid.UUID, value string) error
	Delete(ctx context.Context, id uuid.UUID) error
	WithTx(tx transaction.Tx) LoadTestRepository
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
