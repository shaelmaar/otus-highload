package interceptors

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/ctxcarrier"
)

func UnaryZapLogger(l *zap.Logger, serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()

		ctxLog := l.With(serverCallFields(l, info.FullMethod)...)

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			fields := make([]zap.Field, 0, len(md))

			for k, v := range md {
				fields = append(fields, zap.Strings(fmt.Sprintf("grpc.metadata.%s", k), v))
			}

			fields = append(fields, zap.String("x-request-id", ctxcarrier.ExtractRequestID(ctx)))

			ctxLog = ctxLog.With(fields...)
		}

		ctx = ctxcarrier.InjectLogger(ctx, ctxLog)

		resp, err := handler(ctx, req)

		code := status.Code(err)

		if code == codes.Unknown {
			ctxcarrier.ExtractLogger(ctx).Error(
				"request",
				zap.Error(err),
				zap.String("grpc.code", code.String()),
				zap.Duration("duration", time.Since(startTime)),
				zap.String("component", serviceName),
			)

			return resp, err
		}

		ctxcarrier.ExtractLogger(ctx).Info(
			"request",
			zap.Error(err),
			zap.String("grpc.code", code.String()),
			zap.Duration("duration", time.Since(startTime)),
			zap.String("component", serviceName),
		)

		return resp, err
	}
}

func serverCallFields(l *zap.Logger, fullMethodString string) []zapcore.Field {
	defer func() {
		if r := recover(); r != nil {
			l.Error(
				"[PANIC RECOVER] while parse gRPCInfo.FullMethod",
				zap.Stack("stack"),
				zap.Error(fmt.Errorf("%v", r)),
			)
		}
	}()

	pkg := strings.Split(path.Dir(fullMethodString)[1:], ".")[0]
	service := strings.Split(path.Dir(fullMethodString)[1:], ".")[1]
	method := path.Base(fullMethodString)

	return []zapcore.Field{
		zap.String("grpc.package", pkg),
		zap.String("grpc.service", service),
		zap.String("grpc.method", method),
	}
}
