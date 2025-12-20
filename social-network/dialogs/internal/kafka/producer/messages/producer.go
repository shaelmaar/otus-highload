package messages

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/shaelmaar/otus-highload/social-network/dialogs/gen/kafka/producer/proto"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

const TopicName = "dialogs.messages.event.v1"

type Producer struct {
	writer *kafka.Writer
}

func New(writer *kafka.Writer) (*Producer, error) {
	if utils.IsNil(writer) {
		return nil, errors.New("writer is nil")
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) MessageCreated(ctx context.Context, event dto.MessageCreatedEvent) error {
	msg := pb.DialogsMessagesEventV1{
		EventId:   uuid.New().String(),
		Timestamp: timestamppb.Now(),
		DialogId:  event.DialogID,
		Event: &pb.DialogsMessagesEventV1_Created{
			Created: &pb.DialogsMessagesEventV1_MessageCreated{
				MessageId: event.MessageID,
				From:      event.From.String(),
				To:        event.To.String(),
			},
		},
	}

	key := []byte(event.DialogID)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message created event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{Value: data, Key: key})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) MessagesRead(ctx context.Context, event dto.MessagesReadEvent) error {
	msg := pb.DialogsMessagesEventV1{
		EventId:   uuid.New().String(),
		Timestamp: timestamppb.Now(),
		DialogId:  event.DialogID,
		Event: &pb.DialogsMessagesEventV1_MessagesRead_{
			MessagesRead: &pb.DialogsMessagesEventV1_MessagesRead{
				MessageIds: event.MessageIDs,
				To:         event.To.String(),
			},
		},
	}

	key := []byte(event.DialogID)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshal messages read event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{Value: data, Key: key})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
