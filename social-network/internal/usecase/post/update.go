package post

import (
	"context"
	"fmt"
	"time"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

func (u *UseCases) Update(ctx context.Context, input dto.PostUpdate) error {
	f := func(ctx context.Context, tx transaction.Tx) error {
		repo := u.repo.WithTx(tx)

		post, err := repo.PostLockByID(ctx, input.ID)
		if err != nil {
			return fmt.Errorf("failed to lock post by id: %w", err)
		}

		if post.AuthorUserID != input.UserID {
			return domain.ErrPostNotFoundForUser
		}

		post.UpdatedAt = time.Now()
		post.Content = input.Content

		err = repo.Update(ctx, post)
		if err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}

		return nil
	}

	err := u.tx.Exec(ctx, f, nil)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return nil
}
