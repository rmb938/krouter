package v0

import (
	"math"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v0"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("find-coordinator-v2-handler")

	request := message.(*v0.Request)

	response := &v0.Response{}

	kafkaResponse, err := broker.GetController().FindGroupCoordinator(request.Key)
	if err != nil {
		return nil, err
	}

	if kafkaResponse.ErrorCode != int16(errors.None) {
		response.ErrCode = errors.KafkaError(kafkaResponse.ErrorCode)
	} else {
		response.ErrCode = errors.None
		response.NodeID = math.MaxInt32
		response.Host = broker.AdvertiseListener.IP.String()
		response.Port = int32(broker.AdvertiseListener.Port)
	}

	return response, nil
}
