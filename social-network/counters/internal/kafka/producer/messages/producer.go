package messages

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/shaelmaar/otus-highload/social-network/counters/gen/kafka/producer/proto"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/counters/pkg/utils"
)

const TopicName = "counters.messages.event.v1"

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

func (p *Producer) UnreadMessagesIncremented(ctx context.Context, event dto.UnreadDialogMessagesIncremented) error {
	msg := pb.CountersMessagesEventV1{
		EventId:   uuid.New().String(),
		Timestamp: timestamppb.Now(),
		DialogId:  event.DialogID,
		Event: &pb.CountersMessagesEventV1_UnreadIncremented_{
			UnreadIncremented: &pb.CountersMessagesEventV1_UnreadIncremented{
				MessageId: event.MessageID,
			},
		},
	}

	if !event.Success {
		msg.Event = &pb.CountersMessagesEventV1_UnreadIncrementFailed_{
			UnreadIncrementFailed: &pb.CountersMessagesEventV1_UnreadIncrementFailed{
				MessageId: event.MessageID,
			},
		}
	}

	key := []byte(event.DialogID)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshal unread incremented event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{Value: data, Key: key})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) UnreadMessagesDecremented(ctx context.Context, event dto.UnreadDialogMessagesDecremented) error {
	msg := pb.CountersMessagesEventV1{
		EventId:   uuid.New().String(),
		Timestamp: timestamppb.Now(),
		DialogId:  event.DialogID,
		Event: &pb.CountersMessagesEventV1_UnreadDecremented_{
			UnreadDecremented: &pb.CountersMessagesEventV1_UnreadDecremented{
				MessageIds: event.MessageIDs,
			},
		},
	}

	if !event.Success {
		msg.Event = &pb.CountersMessagesEventV1_UnreadDecrementFailed_{
			UnreadDecrementFailed: &pb.CountersMessagesEventV1_UnreadDecrementFailed{
				MessageIds: event.MessageIDs,
			},
		}
	}

	key := []byte(event.DialogID)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshal unread decremented event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{Value: data, Key: key})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
