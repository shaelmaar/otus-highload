package user

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (uc *UseCases) Login(ctx context.Context, dto dto.Login) (string, error) {
	var out string

	user, err := uc.repo.GetByID(ctx, dto.UserID)
	if err != nil {
		return out, fmt.Errorf("failed to get user by id: %w", err)
	}

	if !utils.CheckPasswordHash(dto.Password, user.PasswordHash) {
		return out, domain.ErrInvalidCredentials
	}

	out, err = uc.auth.GenerateToken(user.ID.String())
	if err != nil {
		return out, fmt.Errorf("failed to generate token: %w", err)
	}

	return out, nil
}
