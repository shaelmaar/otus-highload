package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (uc *UseCases) Register(ctx context.Context, dto dto.RegisterDTO) (domain.User, error) {
	var out domain.User

	passwordHash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return out, fmt.Errorf("failed to hash password: %w", err)
	}

	out = domain.User{
		ID:           uuid.New(),
		PasswordHash: passwordHash,
		FirstName:    dto.Name,
		SecondName:   dto.SecondName,
		BirthDate:    dto.BirthDate,
		Gender:       domain.GenderMale,
		Biography:    dto.Biography,
		City:         dto.City,
	}

	err = uc.repo.Create(ctx, out)
	if err != nil {
		return out, fmt.Errorf("failed to create user: %w", err)
	}

	return out, nil
}
