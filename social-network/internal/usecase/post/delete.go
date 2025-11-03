package post

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

func (u *UseCases) Delete(ctx context.Context, input dto.PostDelete) error {
	f := func(ctx context.Context, tx transaction.Tx) error {
		repo := u.repo.WithTx(tx)

		post, err := repo.PostLockByID(ctx, input.ID)
		if err != nil {
			return fmt.Errorf("failed to lock post by id: %w", err)
		}

		if post.AuthorUserID != input.UserID {
			return domain.ErrPostNotFoundForUser
		}

		err = repo.Delete(ctx, input.ID)
		if err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}

		return nil
	}

	err := u.tx.Exec(ctx, f, nil)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	friendIDs, err := u.friendRepo.Slave().GetFriendUserIDs(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to get friend ids: %w", err)
	}

	err = u.publishUserFeedChunkedTasks(ctx, friendIDs)
	if err != nil {
		return fmt.Errorf("failed to publish user feed chunked tasks: %w", err)
	}

	return nil
}
