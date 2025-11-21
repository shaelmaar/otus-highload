package ctxcarrier

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

func InjectLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

func ExtractLogger(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKey{}).(*zap.Logger)

	if !ok {
		return zap.L()
	}

	return l
}
