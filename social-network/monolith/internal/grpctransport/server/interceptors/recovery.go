package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/shaelmaar/otus-highload/social-network/internal/ctxcarrier"
)

func UnaryRecover(
	ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	panicked := true

	defer func() {
		if r := recover(); r != nil || panicked {
			err = status.Errorf(codes.Unknown, "panic recover: %v", r)

			ctxcarrier.ExtractLogger(ctx).Error("[PANIC RECOVER]", zap.Stack("stack"), zap.Error(err))
		}
	}()

	resp, err = handler(ctx, req)
	panicked = false

	return resp, err
}
