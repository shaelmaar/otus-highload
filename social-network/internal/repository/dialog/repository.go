package dialog

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/datetime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Repository struct {
	mongoDB       *mongo.Database
	tarantoolConn *tarantool.Connection
}

func New(
	mongoDB *mongo.Database,
	tarantoolConn *tarantool.Connection,
) (*Repository, error) {
	if utils.IsNil(mongoDB) {
		return nil, errors.New("mongoDB is nil")
	}

	if utils.IsNil(tarantoolConn) {
		return nil, errors.New("tarantool connection is nil")
	}

	return &Repository{
		mongoDB:       mongoDB,
		tarantoolConn: tarantoolConn,
	}, nil
}

func (r *Repository) CreateDialogMessage(ctx context.Context, message domain.DialogMessage) error {
	collection := r.mongoDB.Collection("dialogMessages")

	_, err := collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to insert new dialog message: %w", err)
	}

	return nil
}

func (r *Repository) CreateDialogMessageTarantool(_ context.Context, message domain.DialogMessage) error {
	dt, err := datetime.NewDatetime(message.CreatedAt.UTC())
	if err != nil {
		return fmt.Errorf("failed to create datetime: %w", err)
	}

	_, err = r.tarantoolConn.Call("add_message",
		[]any{
			message.DialogID.Hex(), message.From, message.To, message.Text, dt,
		})
	if err != nil {
		return fmt.Errorf("failed to add message in tarantool: %w", err)
	}

	return nil
}

func (r *Repository) GetMessagesByDialog(
	ctx context.Context, dialogID primitive.ObjectID) ([]domain.DialogMessage, error) {
	cursor, err := r.mongoDB.Collection("dialogMessages").Find(ctx, bson.M{
		"dialogID": dialogID,
	}, options.Find().SetSort(bson.M{"createdAt": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find dialog messages: %w", err)
	}
	defer cursor.Close(ctx) //nolint:errcheck // можно не проверять здесь ошибку.

	var messages []domain.DialogMessage

	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dialog messages: %w", err)
	}

	return messages, nil
}

func (r *Repository) GetMessagesByDialogTarantool(
	_ context.Context, dialogID primitive.ObjectID) ([]domain.DialogMessage, error) {
	resp, err := r.tarantoolConn.Call("get_messages_by_dialog", []any{dialogID.Hex()})
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
