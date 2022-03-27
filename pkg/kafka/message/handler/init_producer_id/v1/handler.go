package v1

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	request := message.(*v1.Request)

	response := &v1.Response{}

	kafkaResponse, err := broker.GetController().InitProducer(request.TransactionalID, request.TransactionTimeoutDuration)
	if err != nil {
		return nil, err
	}

	response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond
	response.ErrCode = errors.KafkaError(kafkaResponse.ErrorCode)

	response.ProducerID = kafkaResponse.ProducerID
	response.ProducerEpoch = kafkaResponse.ProducerEpoch

	return response, nil
}
