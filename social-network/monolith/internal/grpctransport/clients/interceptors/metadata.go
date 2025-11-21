package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/shaelmaar/otus-highload/social-network/internal/ctxcarrier"
)

const requestIDKey = "x-request-id"

func UnaryClientMetadata(source string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		callOpts ...grpc.CallOption,
	) error {
		requestID := ctxcarrier.ExtractRequestID(ctx)
		if requestID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, requestIDKey, requestID)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "client_name", source)

		return invoker(ctx, method, req, reply, cc, callOpts...)
	}
}
