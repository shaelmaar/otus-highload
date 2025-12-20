package countersmessages

import (
	"context"
)

type DialogMessagesUseCases interface {
	MarkMessageAsSent(ctx context.Context, dialogID string, messageID uint64) error
	MarkMessageAsFailed(ctx context.Context, dialogID string, messageID uint64) error
	MarkMessagesAsRead(ctx context.Context, dialogID string, messageIDs []uint64) error
	MarkMessagesAsSentAfterReading(ctx context.Context, dialogID string, messageIDs []uint64) error
}
