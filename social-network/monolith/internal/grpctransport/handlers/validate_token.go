package handlers

import (
	"context"
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/gen/servergrpc"
	grpcServer "github.com/shaelmaar/otus-highload/social-network/internal/grpctransport/server"
	"github.com/shaelmaar/otus-highload/social-network/internal/service/auth"
)

func (h *Handlers) ValidateToken(
	_ context.Context, req *servergrpc.ValidateTokenRequest) (*servergrpc.ValidateTokenReply, error) {
	if req.Token == "" {
		return nil, grpcServer.GRPCValidationError(servergrpc.ValidateTokenReply_VALIDATION_ERROR,
			errors.New("token is required"))
	}

	userID, err := h.auth.ValidateToken(req.Token)

	switch {
	case errors.Is(err, auth.ErrTokenInvalid):
		return nil, grpcServer.GRPCBusinessError(servergrpc.ValidateTokenReply_TOKEN_INVALID, err)
	case errors.Is(err, auth.ErrTokenExpired):
		return nil, grpcServer.GRPCBusinessError(servergrpc.ValidateTokenReply_TOKEN_EXPIRED, err)
	case err != nil:
		return nil, grpcServer.GRPCUnknownError(servergrpc.ValidateTokenReply_UNSPECIFIED, err)
	}

	return &servergrpc.ValidateTokenReply{
		UserId: userID,
	}, nil
}
