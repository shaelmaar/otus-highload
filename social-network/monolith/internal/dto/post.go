package dto

import (
	"fmt"

	"github.com/google/uuid"
)

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

type UserFeedUpdateTask struct {
	UserID uuid.UUID `json:"user_id"`
}

func (t UserFeedUpdateTask) Info() string {
	return t.UserID.String()
}

type UserFeedChunkedUpdateTask struct {
	UserIDs []uuid.UUID `json:"user_ids"`
}

func (t UserFeedChunkedUpdateTask) Info() string {
	switch {
	case len(t.UserIDs) == 0:
		return ""
	case len(t.UserIDs) == 1:
		return t.UserIDs[0].String()
	default:
		return fmt.Sprintf("%s-%s:%d", t.UserIDs[0], t.UserIDs[len(t.UserIDs)-1], len(t.UserIDs))
	}
}

type PostCreatedChunkedTask struct {
	UserIDs  []uuid.UUID `json:"user_ids"`
	PostID   uuid.UUID   `json:"post_id"`
	Text     string      `json:"text"`
	AuthorID uuid.UUID   `json:"author_id"`
}

func (t PostCreatedChunkedTask) Info() string {
	switch {
	case len(t.UserIDs) == 0:
		return ""
	case len(t.UserIDs) == 1:
		return t.UserIDs[0].String()
	default:
		return fmt.Sprintf("%s-%s:%d", t.UserIDs[0], t.UserIDs[len(t.UserIDs)-1], len(t.UserIDs))
	}
}
