package dialogs

import (
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	dialogsGRPC "github.com/shaelmaar/otus-highload/social-network/gen/clientgrpc/dialogs"
	"github.com/shaelmaar/otus-highload/social-network/internal/grpctransport/clients/interceptors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type NewClientOptions struct {
	GRPCAddr string
	TLS      bool
	Timeout  *time.Duration

	UnaryInterceptors []grpc.UnaryClientInterceptor `exhaustruct:"optional"`

	DialOptions []grpc.DialOption
}

type Client struct {
	grpcConn *grpc.ClientConn

	DialogsService dialogsGRPC.DialogsServiceV1Client
}

type GRPCErrors interface {
	dialogsGRPC.CreateMessageReply_ERR_INFO_REASON | dialogsGRPC.GetDialogMessagesReply_ERR_INFO_REASON
	String() string
}

func IsErr[T GRPCErrors](err error, grpcErr ...T) bool {
	if err == nil || len(grpcErr) < 1 {
		return false
	}

	statusErr, ok := status.FromError(err)
	if !ok {
		return false
	}

	for _, detail := range statusErr.Details() {
		reason, ok := detail.(*errdetails.ErrorInfo)
		if !ok {
			continue
		}

		for _, e := range grpcErr {
			if reason.Reason == e.String() {
				return true
			}
		}
	}

	return false
}

func NewGRPCClient(opts *NewClientOptions) (*Client, error) {
	var creds credentials.TransportCredentials

	if opts.TLS {
		creds = credentials.NewTLS(&tls.Config{}) //nolint:gosec,exhaustruct // берем настройки из OS.
	} else {
		creds = insecure.NewCredentials()
	}

	unaryInter := []grpc.UnaryClientInterceptor{
		interceptors.UnaryClientMetadata("social-network"),
		interceptors.UnaryClientTimeout(opts.Timeout),
	}
	unaryInter = append(unaryInter, opts.UnaryInterceptors...)

	dOpts := make([]grpc.DialOption, 0, len(opts.DialOptions)+3)

	dOpts = append(
		dOpts,
		grpc.WithTransportCredentials(creds),
		grpc.WithChainUnaryInterceptor(unaryInter...),
	)

	dOpts = append(dOpts, opts.DialOptions...)

	clientConn, err := grpc.NewClient(opts.GRPCAddr, dOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to init new grpc client: %w", err)
	}

	client := new(Client)

	client.grpcConn = clientConn
	client.DialogsService = dialogsGRPC.NewDialogsServiceV1Client(clientConn)

	return client, nil
}
