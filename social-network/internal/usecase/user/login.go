package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (uc *UseCases) Login(ctx context.Context, dto dto.LoginDTO) (domain.UserToken, error) {
	var out domain.UserToken

	user, err := uc.repo.GetByID(ctx, dto.UserID)
	if err != nil {
		return out, fmt.Errorf("failed to get user by id: %w", err)
	}

	if !utils.CheckPasswordHash(dto.Password, user.PasswordHash) {
		return out, domain.ErrInvalidCredentials
	}

	f := func(ctx context.Context, tx transaction.Tx) error {
		err := uc.repo.WithTx(tx).DeleteUserTokens(ctx, dto.UserID)
		if err != nil {
			return fmt.Errorf("failed to delete user tokens: %w", err)
		}

		out = domain.UserToken{
			ID:        0,
			UserID:    dto.UserID,
			Token:     uuid.New().String(),
			ExpiresAt: time.Now().Add(30 * time.Minute),
			CreatedAt: time.Time{},
		}

		out.ID, err = uc.repo.CreateUserToken(ctx, out)
		if err != nil {
			return fmt.Errorf("failed to create user token: %w", err)
		}

		return nil
	}

	err = uc.tx.Exec(ctx, f, nil)
	if err != nil {
		return out, fmt.Errorf("failed to handle login transaction: %w", err)
	}

	return out, nil
}
