package dialog

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/datetime"

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

func (r *Repository) CreateDialogMessage(_ context.Context, message domain.DialogMessage) error {
	dt, err := datetime.NewDatetime(message.CreatedAt.UTC())
	if err != nil {
		return fmt.Errorf("failed to create datetime: %w", err)
	}

	_, err = r.tarantoolConn.Call("add_message",
		[]any{
			message.DialogID, message.From, message.To, message.Text, dt,
		})
	if err != nil {
		return fmt.Errorf("failed to add message in tarantool: %w", err)
	}

	return nil
}

func (r *Repository) GetMessagesByDialog(
	_ context.Context, dialogID string) ([]domain.DialogMessage, error) {
	resp, err := r.tarantoolConn.Call("get_messages_by_dialog", []any{dialogID})
	if err != nil {
		return nil, fmt.Errorf("failed to find dialog messages in tarantool: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, nil
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

		from, _ := t["from"].(uuid.UUID)
		to, _ := t["to"].(uuid.UUID)
		text, _ := t["text"].(string)
		createdAt, _ := t["created_at"].(datetime.Datetime)

		msg := domain.DialogMessage{
			From:      from,
			To:        to,
			DialogID:  dialogID,
			Text:      text,
			CreatedAt: createdAt.ToTime(),
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
