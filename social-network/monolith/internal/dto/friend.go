package dto

import "github.com/google/uuid"

type FriendSet struct {
	UserID   uuid.UUID
	FriendID uuid.UUID
}

type FriendDelete struct {
	UserID   uuid.UUID
	FriendID uuid.UUID
}
