package dialog

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

type UseCases struct {
	repo          domain.DialogRepository
	kafkaProducer KafkaProducer
}

func New(repo domain.DialogRepository, kafkaProducer KafkaProducer) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repository is nil")
	}

	if utils.IsNil(kafkaProducer) {
		return nil, errors.New("kafka producer is nil")
	}

	return &UseCases{
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}, nil
}
