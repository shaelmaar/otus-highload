package user

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (uc *UseCases) Search(ctx context.Context, dto dto.SearchDTO) ([]domain.User, error) {
	return uc.repo.Slave().GetByFirstNameLastName(ctx, dto.FirstName, dto.LastName)
}
