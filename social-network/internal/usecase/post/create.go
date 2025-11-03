package post

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (u *UseCases) Create(ctx context.Context, input dto.PostCreate) (uuid.UUID, error) {
	post := domain.Post{
		ID:           uuid.New(),
		Content:      input.Content,
		AuthorUserID: input.UserID,
		CreatedAt:    time.Now(),

		// не используется.
		UpdatedAt: time.Now(),
	}

	err := u.repo.Create(ctx, post)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create post: %w", err)
	}

	friendIDs, err := u.friendRepo.Slave().GetFriendUserIDs(ctx, input.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get friend ids: %w", err)
	}

	err = u.publishUserFeedChunkedTasks(ctx, friendIDs)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to publish user feed chunked tasks: %w", err)
	}

	err = u.publishPostCreatedChunkedTasks(ctx, friendIDs, post.ID, post.AuthorUserID, post.Content)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to publish post created chunked tasks: %w", err)
	}

	return post.ID, nil
}

func (u *UseCases) publishUserFeedChunkedTasks(ctx context.Context, friendIDs []uuid.UUID) error {
	for chunk := range utils.ChunkSlice(ctx, friendIDs, userFeedChunkSize) {
		err := u.userFeedChunkedProducer.Publish(ctx, dto.UserFeedChunkedUpdateTask{UserIDs: chunk})
		if err != nil {
			return fmt.Errorf("failed to publish user feed chunked update: %w", err)
		}
	}

	return nil
}

func (u *UseCases) publishPostCreatedChunkedTasks(
	ctx context.Context, friendIDs []uuid.UUID,
	postID, authorUserID uuid.UUID,
	postText string,
) error {
	for chunk := range utils.ChunkSlice(ctx, friendIDs, userFeedChunkSize) {
		err := u.postCreatedProducer.Publish(ctx, dto.PostCreatedChunkedTask{
			UserIDs:  chunk,
			PostID:   postID,
			Text:     postText,
			AuthorID: authorUserID,
		})
		if err != nil {
			return fmt.Errorf("failed to publish user feed chunked update: %w", err)
		}
	}

	return nil
}
