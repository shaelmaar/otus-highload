package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

const (
	defaultWorkerCount   = 2
	defaultPrefetchCount = 2
	defaultWorkerTimeout = 3 * time.Second
)

type Consumer[T Message] struct {
	queueName     string
	handler       func(context.Context, T) error
	logger        *zap.Logger
	conn          *amqp.Connection
	ch            *amqp.Channel
	workerCount   int
	prefetchCount int
	workerTimeout time.Duration
	messageTTL    time.Duration
}

func NewConsumer[T Message](
	url, queueName, exchangeName string,
	handler func(context.Context, T) error,
	logger *zap.Logger,
	opts ...ConsumerOption[T],
) (*Consumer[T], error) {
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

	c := &Consumer[T]{
		queueName:     queueName,
		handler:       handler,
		workerCount:   defaultWorkerCount,
		prefetchCount: defaultPrefetchCount,
		workerTimeout: defaultWorkerTimeout,
		messageTTL:    defaultMessageTTL,
		logger:        logger,
		conn:          conn,
		ch:            ch,
	}

	for _, opt := range opts {
		opt(c)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": c.messageTTL.Milliseconds(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	if exchangeName != "" {
		err = ch.ExchangeDeclare(exchangeName, "fanout", false, false, false, false, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to declare an exchange: %w", err)
		}

		err = ch.QueueBind(queue.Name, "", exchangeName, false, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to bind a queue: %w", err)
		}

		c.queueName = queue.Name
	}

	err = ch.Qos(
		c.prefetchCount,
		0,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return c, nil
}

//nolint:errcheck,gosec // проверять при закрытии необязательно.
func (c *Consumer[T]) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer[T]) Consume(ctx context.Context) error {
	// Получаем сообщения из очереди
	messages, err := c.ch.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	var wg sync.WaitGroup
	messageChan := make(chan amqp.Delivery, c.workerCount)

	for i := 0; i < c.workerCount; i++ {
		wg.Add(1)
		go c.worker(ctx, i, &wg, messageChan)
	}

	c.logger.Info("Started workers for queue",
		zap.String("queue_name", c.queueName),
		zap.Int("worker_count", c.workerCount))

	for {
		select {
		case <-ctx.Done():
			close(messageChan)

			wg.Wait()

			return nil
		case message, ok := <-messages:
			if !ok {
				close(messageChan)
				wg.Wait()

				return nil
			}

			messageChan <- message
		}
	}
}

func (c *Consumer[T]) worker(ctx context.Context, n int, wg *sync.WaitGroup, messages <-chan amqp.Delivery) {
	defer wg.Done()

	logger := c.logger.With(zap.String("queue_name", c.queueName),
		zap.Int("worker_num", n))
	logger.Info("started worker for queue")

	for message := range messages {
		var messagePayload T

		err := json.Unmarshal(message.Body, &messagePayload)
		if err != nil {
			logger.Error("failed to unmarshal message", zap.Error(err))

			continue
		}

		func() {
			workerCtx, workerCancel := context.WithTimeout(ctx, c.workerTimeout)
			defer workerCancel()

			if err := c.handler(workerCtx, messagePayload); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}

				logger.Error("failed to handle message", zap.Error(err))

				err = message.Nack(false, true)
				if err != nil {
					logger.Error("failed to nack message", zap.Error(err))
				}

				return
			}

			err = message.Ack(false)
			if err != nil {
				logger.Error("failed to ack message", zap.Error(err))
			}

			logger.Info("message acked", zap.String("queue_name", c.queueName))
		}()
	}

	logger.Info("worker stopped")
}
