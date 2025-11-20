package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/shaelmaar/otus-highload/social-network/gen/servergrpc"
	"github.com/shaelmaar/otus-highload/social-network/internal/grpctransport/server/interceptors"
)

type GRPCHandlers interface {
	ValidateToken(ctx context.Context, req *servergrpc.ValidateTokenRequest) (*servergrpc.ValidateTokenReply, error)
}

type Server struct {
	Handlers  GRPCHandlers
	Validator *validator.Validate

	logger *zap.Logger

	grpcServer *grpc.Server
}

func (s *Server) Serve(port *int) error {
	if port == nil {
		defaultPort := 8002
		port = &defaultPort
	}

	//nolint:noctx // контекста нет.
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		return err
	}

	s.logger.Info("starting gRPC server", zap.Int("port", *port))

	if err := s.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("error while Serve grpc | %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

type NewServerOptions struct {
	Logger       *zap.Logger
	GRPCHandlers GRPCHandlers
	Validator    *validator.Validate

	UnaryInterceptors []grpc.UnaryServerInterceptor

	ServerOptions []grpc.ServerOption
}

func New(opts *NewServerOptions) (*Server, error) {
	unaryInter := unaryInterceptors(opts.Logger, "social-network")
	unaryInter = append(unaryInter, opts.UnaryInterceptors...)

	srvOpts := make([]grpc.ServerOption, 0, len(opts.ServerOptions)+3)

	srvOpts = append(
		srvOpts,
		grpc.ChainUnaryInterceptor(unaryInter...),
	)

	srvOpts = append(srvOpts, opts.ServerOptions...)

	grpcServer := grpc.NewServer(srvOpts...)

	if opts.Validator == nil {
		opts.Validator = validator.New()
	}

	s := Server{
		grpcServer: grpcServer,
		Validator:  opts.Validator,
		logger:     opts.Logger,
		Handlers:   opts.GRPCHandlers,
	}

	//nolint:exhaustruct // UnsafeAuthServiceV1Server эмбеддинг интерфейса.
	servergrpc.RegisterAuthServiceV1Server(grpcServer, &gRPCAuthServiceService{Server: s})

	reflection.Register(grpcServer)

	return &s, nil
}

type gRPCAuthServiceService struct {
	servergrpc.UnsafeAuthServiceV1Server
	Server
}

func (s *gRPCAuthServiceService) ValidateToken(
	ctx context.Context, req *servergrpc.ValidateTokenRequest) (*servergrpc.ValidateTokenReply, error) {
	return s.Handlers.ValidateToken(ctx, req)
}

type GRPCErrors interface {
	servergrpc.ValidateTokenReply_ERR_INFO_REASON
	String() string
}

func GRPCValidationError[T GRPCErrors](reason T, err error) error {
	return gRPCError(codes.InvalidArgument, reason, err)
}

func GRPCBusinessError[T GRPCErrors](reason T, err error) error {
	return gRPCError(codes.FailedPrecondition, reason, err)
}

func GRPCUnknownError[T GRPCErrors](reason T, err error) error {
	return gRPCError(codes.Unknown, reason, err)
}

func GRPCCustomError[T GRPCErrors](code codes.Code, reason T, err error) error {
	return gRPCError(code, reason, err)
}

func gRPCError[T GRPCErrors](code codes.Code, reason T, serviceErr error) error {
	if serviceErr == nil {
		serviceErr = errors.New("error not set")
	}

	st, err := status.New(code, serviceErr.Error()).WithDetails(
		&errdetails.ErrorInfo{
			Reason:   reason.String(),
			Domain:   "social-network",
			Metadata: nil,
		},
	)
	if err != nil {
		panic(fmt.Sprintf("unexpected error attaching metadata: %v", err))
	}

	var vErrors validator.ValidationErrors
	if errors.As(serviceErr, &vErrors) {
		var br errdetails.BadRequest

		for _, vErr := range vErrors {
			//nolint:exhaustruct // остальное неважно.
			v := &errdetails.BadRequest_FieldViolation{
				Field:       vErr.StructNamespace(),
				Description: vErr.Tag(),
			}

			br.FieldViolations = append(br.FieldViolations, v)
		}

		st, err = st.WithDetails(&br)
		if err != nil {
			panic(fmt.Sprintf("unexpected error attaching metadata: %v", err))
		}
	}

	return st.Err()
}

func unaryInterceptors(l *zap.Logger, serviceName string) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		interceptors.UnaryRequestID(l),
		interceptors.UnaryZapLogger(l, serviceName),
		interceptors.UnaryRecover,
	}
}
