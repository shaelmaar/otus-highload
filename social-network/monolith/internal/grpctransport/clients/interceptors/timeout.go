package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func UnaryClientTimeout(timeout *time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		callOpts ...grpc.CallOption,
	) error {
		if timeout == nil {
			return invoker(ctx, method, req, reply, cc, callOpts...)
		}

		ctx, cancel := context.WithTimeout(ctx, *timeout)

		defer cancel()

		return invoker(ctx, method, req, reply, cc, callOpts...)
	}
}
