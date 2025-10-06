package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

const (
	defaultMessageTTL = 15 * time.Minute
)

type Message interface {
	Info() string
}

type Producer[T Message] struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	queue  amqp.Queue
	logger *zap.Logger

	messageTTL time.Duration
}

func NewProducer[T Message](
	url, queueName string,
	logger *zap.Logger,
	opts ...ProducerOption[T],
) (*Producer[T], error) {
	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	p := &Producer[T]{
		conn:       conn,
		ch:         ch,
		queue:      amqp.Queue{}, //nolint:exhaustruct // пустая структура, определяется ниже.
		logger:     logger,
		messageTTL: defaultMessageTTL,
	}

	for _, opt := range opts {
		opt(p)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": p.messageTTL.Milliseconds(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	p.queue = queue

	return p, nil
}

//nolint:errcheck,gosec // проверять при закрытии необязательно.
func (p *Producer[T]) Close() {
	if p.ch != nil {
		p.ch.Close()
	}

	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Producer[T]) Publish(ctx context.Context, message T) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	err = p.ch.PublishWithContext(ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{ //nolint:exhaustruct // остальное пока не нужно.
			ContentType: "application/json",
			Body:        payload,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	p.logger.Info("message published",
		zap.String("queue_name", p.queue.Name),
		zap.String("message_info", message.Info()),
	)

	return nil
}
