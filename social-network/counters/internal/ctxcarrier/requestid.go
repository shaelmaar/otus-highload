package ctxcarrier

import (
	"context"
)

type requestIDKey struct{}

func InjectRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

func ExtractRequestID(ctx context.Context) string {
	value := ctx.Value(requestIDKey{})
	if str, ok := value.(string); ok {
		return str
	}

	return ""
}
