package dialogsmessages

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	pb "github.com/shaelmaar/otus-highload/social-network/counters/gen/kafka/consumer/proto"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/counters/pkg/utils"
)

const TopicName = "dialogs.messages.event.v1"

var errInvalidMessage = errors.New("invalid message")

type Option func(*Consumer)

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(listener *Consumer) {
		listener.shutdownTimeout = timeout
	}
}

type Consumer struct {
	reader   *kafka.Reader
	useCases DialogMessagesCounterUseCases
	logger   *zap.Logger

	ctx             context.Context
	cancel          context.CancelFunc
	shutdownTimeout time.Duration
	jobCtx          context.Context
	jobCancel       context.CancelFunc
}

func New(
	reader *kafka.Reader,
	useCases DialogMessagesCounterUseCases,
	logger *zap.Logger,
	opts ...Option,
) (*Consumer, error) {
	if utils.IsNil(reader) {
		return nil, errors.New("reader is nil")
	}

	if utils.IsNil(useCases) {
		return nil, errors.New("use cases is nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := Consumer{
		reader:   reader,
		useCases: useCases,
		logger:   logger,

		ctx:             ctx,
		cancel:          cancel,
		shutdownTimeout: 10 * time.Second,
		jobCtx:          nil,
		jobCancel:       nil,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c, nil
}

func (c *Consumer) Consume() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			func() {
				c.jobCtx, c.jobCancel = context.WithCancel(context.Background())
				defer c.jobCancel()

				if err := c.readAndReact(c.jobCtx); err != nil {
					if errors.Is(err, context.Canceled) {
						// c.Close() is called
						return
					}

					c.logger.Error("failed to read and react", zap.Error(err))
				}
			}()
		}
	}
}

func (c *Consumer) Close() error {
	if c.shutdownTimeout > 0 {
		select {
		case <-time.After(c.shutdownTimeout):
			c.jobCancel()
		case <-c.jobCtx.Done():
		}
	}

	c.cancel()

	return nil
}

func (c *Consumer) readAndReact(ctx context.Context) error {
	message, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch message from kafka: %w", err)
	}

	msg := new(pb.DialogsMessagesEventV1)
	if err = proto.Unmarshal(message.Value, msg); err != nil {
		c.logger.Error("failed to unmarshal dialogs messages event", zap.Error(err))

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			return fmt.Errorf("failed to commit message: %w", err)
		}

		return nil
	}

	err = c.doReact(ctx, msg)

	switch {
	case errors.Is(err, errInvalidMessage):
		c.logger.Error("invalid message", zap.Error(err))
	case err != nil:
		return fmt.Errorf("failed to do react: %w", err)
	}

	if err := c.reader.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to commit message: %w", err)
	}

	c.logger.Info("processed message", zap.String("event_id", msg.EventId))

	return nil
}

func (c *Consumer) doReact(ctx context.Context, msg *pb.DialogsMessagesEventV1) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	switch e := msg.Event.(type) {
	case *pb.DialogsMessagesEventV1_Created:
		recipientID, err := uuid.Parse(e.Created.To)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidMessage, err.Error())
		}

		err = c.useCases.IncrementUnreadMessages(ctx, dto.UnreadDialogMessagesIncrement{
			DialogID:       msg.DialogId,
			MessageID:      e.Created.MessageId,
			RecipientID:    recipientID,
			IdempotencyKey: msg.EventId,
		})
		if err != nil {
			return fmt.Errorf("failed to increment unread messages count: %w", err)
		}
	case *pb.DialogsMessagesEventV1_MessagesRead_:
		recipientID, err := uuid.Parse(e.MessagesRead.To)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidMessage, err.Error())
		}

		err = c.useCases.DecrementUnreadMessages(ctx, dto.UnreadDialogMessagesDecrement{
			DialogID:       msg.DialogId,
			MessageIDs:     e.MessagesRead.MessageIds,
			RecipientID:    recipientID,
			IdempotencyKey: msg.EventId,
		})
		if err != nil {
			return fmt.Errorf("failed to decrement unread messages count: %w", err)
		}
	}

	return nil
}
