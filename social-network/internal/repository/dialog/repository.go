package dialog

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
)

type Repository struct {
	db *mongo.Database
}

func New(db *mongo.Database) (*Repository, error) {
	return &Repository{db: db}, nil
}

func (r *Repository) CreateDialogMessage(ctx context.Context, message domain.DialogMessage) error {
	collection := r.db.Collection("dialogMessages")

	_, err := collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to insert new dialog message: %w", err)
	}

	return nil
}

func (r *Repository) GetMessagesByDialog(
	ctx context.Context, dialogID primitive.ObjectID) ([]domain.DialogMessage, error) {
	cursor, err := r.db.Collection("dialogMessages").Find(ctx, bson.M{
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
