package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginDTO struct {
	UserID   uuid.UUID
	Password string
}

type RegisterDTO struct {
	Password   string
	Name       string
	SecondName string
	BirthDate  time.Time
	Biography  string
	City       string
}
