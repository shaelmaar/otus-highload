package auth

import (
	"context"

	"github.com/google/uuid"
)

type userIDKey struct{}

func InjectUserIDToContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func ExtractUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)

	return userID, ok
}
