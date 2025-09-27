package dto

import (
	"time"

	"github.com/google/uuid"
)

type Login struct {
	UserID   uuid.UUID
	Password string
}

type Register struct {
	Password   string
	Name       string
	SecondName string
	BirthDate  time.Time
	Biography  string
	City       string
}

type Search struct {
	FirstName string
	LastName  string
}
