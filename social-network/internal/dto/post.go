package dto

import "github.com/google/uuid"

type PostCreate struct {
	Content string
	UserID  uuid.UUID
}

type PostUpdate struct {
	ID      uuid.UUID
	Content string
	UserID  uuid.UUID
}

type PostDelete struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type GetPostFeed struct {
	UserID uuid.UUID
	Offset int
	Limit  int
}

type GetLastPostsByUserIDs struct {
	UserIDs []uuid.UUID
	Offset  int
	Limit   int
}
