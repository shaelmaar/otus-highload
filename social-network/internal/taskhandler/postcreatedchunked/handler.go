package postcreatedchunked

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/olahol/melody"

	"github.com/shaelmaar/otus-highload/social-network/gen/async"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Handler struct {
	m *melody.Melody
}

func New(m *melody.Melody) (*Handler, error) {
	if utils.IsNil(m) {
		return nil, errors.New("melody is nil")
	}

	return &Handler{
		m: m,
	}, nil
}

func (h *Handler) Handle(_ context.Context, task dto.PostCreatedChunkedTask) error {
	userIDs := utils.SliceToMapAsKeys(task.UserIDs)

	msg := async.PostMessage{
		Payload: async.PostMessagePayload{
			AuthorUserId: utils.Ptr(async.UserIdSchema(task.AuthorID.String())),
			PostId:       utils.Ptr(async.PostIdSchema(task.PostID.String())),
			PostText:     utils.Ptr(async.PostTextSchema(task.Text)),
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal post message: %w", err)
	}

	err = h.m.BroadcastFilter(msgBytes, func(s *melody.Session) bool {
		userID, ok := s.Get("user_id")
		if !ok {
			return false
		}

		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			return false
		}

		_, ok = userIDs[userUUID]

		return ok
	})
	if err != nil {
		return fmt.Errorf("failed to broadcast post created message: %w", err)
	}

	return nil
}
