package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DialogRepository interface {
	CreateDialogMessage(ctx context.Context, message DialogMessage) error
	GetMessagesByDialog(
		ctx context.Context, dialogID primitive.ObjectID) ([]DialogMessage, error)
}

type DialogMessage struct {
	From      uuid.UUID          `bson:"from"`
	To        uuid.UUID          `bson:"to"`
	DialogID  primitive.ObjectID `bson:"dialogID"`
	Text      string             `bson:"text"`
	CreatedAt time.Time          `bson:"createdAt"`
}
