package countersmessages

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	pb "github.com/shaelmaar/otus-highload/social-network/dialogs/gen/kafka/consumer/proto"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

const TopicName = "counters.messages.event.v1"

var errInvalidMessage = errors.New("invalid message")

type Option func(*Consumer)

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(listener *Consumer) {
		listener.shutdownTimeout = timeout
	}
}

type Consumer struct {
	reader   *kafka.Reader
	useCases DialogMessagesUseCases
	logger   *zap.Logger

	ctx             context.Context
	cancel          context.CancelFunc
	shutdownTimeout time.Duration
	jobCtx          context.Context
	jobCancel       context.CancelFunc
}

func New(
	reader *kafka.Reader,
	useCases DialogMessagesUseCases,
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

	msg := new(pb.CountersMessagesEventV1)
	if err = proto.Unmarshal(message.Value, msg); err != nil {
		c.logger.Error("failed to unmarshal counters messages event", zap.Error(err))

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

func (c *Consumer) doReact(ctx context.Context, msg *pb.CountersMessagesEventV1) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	switch e := msg.Event.(type) {
	case *pb.CountersMessagesEventV1_UnreadIncremented_:
		err := c.useCases.MarkMessageAsSent(ctx, msg.DialogId, e.UnreadIncremented.MessageId)
		if err != nil {
			return fmt.Errorf("failed to mark message as failed: %w", err)
		}
	case *pb.CountersMessagesEventV1_UnreadIncrementFailed_:
		err := c.useCases.MarkMessageAsFailed(ctx, msg.DialogId, e.UnreadIncrementFailed.MessageId)
		if err != nil {
			return fmt.Errorf("failed to mark message as failed: %w", err)
		}
	case *pb.CountersMessagesEventV1_UnreadDecremented_:
		err := c.useCases.MarkMessagesAsRead(ctx, msg.DialogId, e.UnreadDecremented.MessageIds)
		if err != nil {
			return fmt.Errorf("failed to mark messages as read: %w", err)
		}
	case *pb.CountersMessagesEventV1_UnreadDecrementFailed_:
		err := c.useCases.MarkMessagesAsSentAfterReading(ctx, msg.DialogId, e.UnreadDecrementFailed.MessageIds)
		if err != nil {
			return fmt.Errorf("failed to mark messages as sent after reading: %w", err)
		}
	default:
		return fmt.Errorf("%w: unknown event '%T'", errInvalidMessage, e)
	}

	return nil
}
