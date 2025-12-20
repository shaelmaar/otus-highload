package dialog

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/counters/pkg/utils"
)

type UseCases struct {
	dialogRepo    domain.DialogRepository
	kafkaProducer KafkaProducer
}

func New(dialogRepo domain.DialogRepository, kafkaProducer KafkaProducer) (*UseCases, error) {
	if utils.IsNil(dialogRepo) {
		return nil, errors.New("dialog repository is nil")
	}

	if utils.IsNil(kafkaProducer) {
		return nil, errors.New("kafka producer is nil")
	}

	return &UseCases{
		dialogRepo:    dialogRepo,
		kafkaProducer: kafkaProducer,
	}, nil
}
