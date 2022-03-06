package v1

import (
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, message message.Message, correlationId int32) error {
	request := message.(*v1.Request)

	response := &v1.Response{}

	response.ThrottleDuration = 0

	response.ErrCode = errors.None
	if request.TransactionalID != nil {
		response.ErrCode = errors.TransactionIDAuthorizationFailed
	}

	response.ErrCode = errors.ClusterAuthorizationFailed // TODO: remove this once we support transactions

	return client.WriteMessage(response, correlationId)
}
