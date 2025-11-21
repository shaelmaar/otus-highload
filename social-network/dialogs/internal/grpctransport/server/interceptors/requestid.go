package interceptors

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/ctxcarrier"
)

const requestIDKey = "x-request-id"

func UnaryRequestID(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		var rid string

		requestID := md.Get(requestIDKey)

		if len(requestID) == 0 {
			rid = uuid.New().String()
		} else {
			rid = requestID[0]
		}

		data := metadata.Pairs(requestIDKey, rid)
		if err := grpc.SetHeader(ctx, data); err != nil {
			l.Error("grpc set header failed", zap.Error(err))
		}

		ctx = ctxcarrier.InjectRequestID(ctx, rid)

		return handler(ctx, req)
	}
}
