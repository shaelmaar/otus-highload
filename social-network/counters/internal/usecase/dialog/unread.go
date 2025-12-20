package dialog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/domain"
)

func (u *UseCases) UnreadMessageCount(ctx context.Context, recipientID, senderID uuid.UUID) (int64, error) {
	res, err := u.dialogRepo.Slave().CountUnreadMessages(ctx, domain.UnreadDialogMessageCountKey{
		RecipientID: recipientID,
		DialogID:    generateDialogID(senderID, recipientID),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}

	return res, nil
}

func generateDialogID(firstID, secondID uuid.UUID) string {
	ids := []string{firstID.String(), secondID.String()}
	sort.Strings(ids)

	hash := sha256.Sum256([]byte(ids[0] + ":" + ids[1]))

	var id [12]byte
	copy(id[:], hash[:12])

	var buf [24]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}
