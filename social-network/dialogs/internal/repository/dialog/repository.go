package dialog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

type Repository struct {
	tarantoolConn *tarantool.Connection
}

func New(
	tarantoolConn *tarantool.Connection,
) (*Repository, error) {
	if utils.IsNil(tarantoolConn) {
		return nil, errors.New("tarantool connection is nil")
	}

	return &Repository{
		tarantoolConn: tarantoolConn,
	}, nil
}

func (r *Repository) CreateDialogMessage(_ context.Context, message domain.DialogMessage) (uint64, error) {
	resp, err := r.tarantoolConn.Call("add_message",
		[]any{
			message.DialogID, message.From, message.To, message.Text, message.State, message.CreatedAt.UnixMicro(),
		})
	if err != nil {
		return 0, fmt.Errorf("failed to add message in tarantool: %w", err)
	}

	if len(resp.Data) == 0 {
		return 0, errors.New("failed to add messaage in tarantool")
	}

	dataSlice, ok := resp.Data[0].([]any)
	if !ok {
		return 0, fmt.Errorf("failed to decode response from tarantool: %w", err)
	}

	data := dataSlice[0].(map[any]any)

	id := data["id"].(uint64)

	return id, nil
}

func (r *Repository) GetMessagesByDialog(
	_ context.Context, dialogID string) ([]domain.DialogMessage, error) {
	resp, err := r.tarantoolConn.Call("get_messages_by_dialog", []any{dialogID})
	if err != nil {
		return nil, fmt.Errorf("failed to find dialog messages in tarantool: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("failed to find dialog messages in tarantool")
	}

	tuples, ok := resp.Data[0].([]any)
	if !ok {
		return nil, fmt.Errorf("failed to decode dialog messages in tarantool: %w", err)
	}

	messages := make([]domain.DialogMessage, 0, len(tuples))

	for _, tuple := range tuples {
		t, ok := tuple.(map[any]any)
		if !ok {
			continue
		}

		id, _ := t["id"].(uint64)
		from, _ := t["from"].(uuid.UUID)
		to, _ := t["to"].(uuid.UUID)
		text, _ := t["text"].(string)
		state, _ := t["state"].(string)
		createdAt, _ := t["created_at"].(uint64)

		msg := domain.DialogMessage{
			ID:        id,
			From:      from,
			To:        to,
			DialogID:  dialogID,
			Text:      text,
			State:     domain.DialogMessageState(state),
			CreatedAt: time.UnixMicro(int64(createdAt)),
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *Repository) MarkMessagesAsReading(
	_ context.Context, dialogID string, readerID uuid.UUID, messageID uint64) ([]uint64, error) {
	resp, err := r.tarantoolConn.Call("mark_messages_as_reading",
		[]any{dialogID, readerID, messageID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to mark messages as reading in tarantool: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("failed to mark messages as reading in tarantool")
	}

	data := resp.Data[0].([]any)

	ids := make([]uint64, 0, len(data))

	for _, idI := range data {
		if id, ok := idI.(uint64); ok {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func (r *Repository) UpdateMessageStateFrom(
	_ context.Context, dialogID string, messageID uint64, fromState, toState domain.DialogMessageState,
) (bool, error) {
	resp, err := r.tarantoolConn.Call("update_message_state_from",
		[]any{dialogID, messageID, string(fromState), string(toState)},
	)
	if err != nil {
		return false, fmt.Errorf("failed to update message state from in tarantool: %w", err)
	}

	if len(resp.Data) == 0 {
		return false, errors.New("failed to update message state from in tarantool")
	}

	dataSlice, ok := resp.Data[0].([]any)
	if !ok {
		return false, fmt.Errorf("failed to decode response from tarantool: %w", err)
	}

	data := dataSlice[0].(map[any]any)

	return data["updated"].(bool), nil
}

func (r *Repository) UpdateMessagesStateFrom(
	_ context.Context, dialogID string, messageIDs []uint64, fromState, toState domain.DialogMessageState,
) ([]uint64, error) {
	resp, err := r.tarantoolConn.Call("update_messages_state_from",
		[]any{dialogID, messageIDs, string(fromState), string(toState)},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update messages state from in tarantool: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("failed to update messages state from in tarantool")
	}

	data := resp.Data[0].([]any)

	ids := make([]uint64, 0, len(data))

	for _, idI := range data {
		if id, ok := idI.(uint64); ok {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
