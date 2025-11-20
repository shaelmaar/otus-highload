package dialog

import (
	"errors"

	dialogsGRPC "github.com/shaelmaar/otus-highload/social-network/gen/clientgrpc/dialogs"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	dialogsClient dialogsGRPC.DialogsServiceV1Client
}

func New(dialogsClient dialogsGRPC.DialogsServiceV1Client) (*UseCases, error) {
	if utils.IsNil(dialogsClient) {
		return nil, errors.New("dialogs client is nil")
	}

	return &UseCases{dialogsClient: dialogsClient}, nil
}
