package dialog

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (u *UseCases) CreateMessage(ctx context.Context, input dto.DialogCreateMessage) error {
	err := u.repo.CreateDialogMessage(ctx, domain.DialogMessage{
		From:      input.From,
		To:        input.To,
		DialogID:  generateDialogID(input.From, input.To),
		Text:      input.Text,
		CreatedAt: input.Time,
	})
	if err != nil {
		return fmt.Errorf("failed to create dialog message: %w", err)
	}

	return nil
}

func generateDialogID(firstID, secondID uuid.UUID) primitive.ObjectID {
	ids := []string{firstID.String(), secondID.String()}
	sort.Strings(ids)

	hash := sha256.Sum256([]byte(ids[0] + ":" + ids[1]))

	var id [12]byte
	copy(id[:], hash[:12])

	return id
}
