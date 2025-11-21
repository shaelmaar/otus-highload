package auth

import (
	"context"
	"errors"
	"fmt"

	monolithGRPC "github.com/shaelmaar/otus-highload/social-network/dialogs/gen/clientgrpc/monolith"
	monolithGRPCClient "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/grpctransport/clients/monolith"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

type Service struct {
	monolithAuthClient monolithGRPC.AuthServiceV1Client
}

func New(
	monolithAuthClient monolithGRPC.AuthServiceV1Client,
) (*Service, error) {
	if utils.IsNil(monolithAuthClient) {
		return nil, errors.New("monolith auth is nil")
	}

	return &Service{
		monolithAuthClient: monolithAuthClient,
	}, nil
}

func (s *Service) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	reply, err := s.monolithAuthClient.ValidateToken(ctx, &monolithGRPC.ValidateTokenRequest{
		Token: tokenString,
	})

	switch {
	case monolithGRPCClient.IsErr(err, monolithGRPC.ValidateTokenReply_TOKEN_INVALID):
		return "", ErrTokenInvalid
	case monolithGRPCClient.IsErr(err, monolithGRPC.ValidateTokenReply_TOKEN_EXPIRED):
		return "", ErrTokenExpired
	case err != nil:
		return "", fmt.Errorf("failed to validate token in monolith auth: %w", err)
	}

	return reply.UserId, nil
}
