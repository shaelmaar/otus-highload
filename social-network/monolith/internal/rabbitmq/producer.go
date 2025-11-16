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
	conn               *amqp.Connection
	ch                 *amqp.Channel
	queueName          string
	fanoutExchangeName string
	logger             *zap.Logger

	messageTTL time.Duration
}

func NewProducer[T Message](
	url, queueName, fanoutExchangeName string,
	logger *zap.Logger,
	opts ...ProducerOption[T],
) (*Producer[T], error) {
	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	if queueName != "" && fanoutExchangeName != "" {
		return nil, errors.New("fanout exchange name and queue name both set")
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
		conn:               conn,
		ch:                 ch,
		queueName:          queueName,
		fanoutExchangeName: fanoutExchangeName,
		logger:             logger,
		messageTTL:         defaultMessageTTL,
	}

	for _, opt := range opts {
		opt(p)
	}

	if queueName != "" {
		_, err = ch.QueueDeclare(
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

		p.queueName = queueName
	}

	if fanoutExchangeName != "" {
		err = ch.ExchangeDeclare(fanoutExchangeName, "fanout", false, false, false, false, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to declare an exchange: %w", err)
		}

		p.fanoutExchangeName = fanoutExchangeName
	}

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
		p.fanoutExchangeName,
		p.queueName,
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
		zap.String("queue_name", p.queueName),
		zap.String("fanout_exchange_name", p.fanoutExchangeName),
		zap.String("message_info", message.Info()),
	)

	return nil
}
