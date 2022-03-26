package v1

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("find-coordinator-v2-handler")

	request := message.(*v2.Request)

	response := &v2.Response{}

	var kafkaResponse *kmsg.FindCoordinatorResponse
	var err error
	if request.KeyType != 0 {
		kafkaResponse, err = broker.GetController().FindTransactionCoordinator(request.Key)
	} else {
		kafkaResponse, err = broker.GetController().FindGroupCoordinator(request.Key)
	}

	if err != nil {
		return nil, err
	}

	response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond

	if kafkaResponse.ErrorCode != int16(errors.None) {
		response.ErrCode = errors.KafkaError(kafkaResponse.ErrorCode)
		response.ErrMessage = kafkaResponse.ErrorMessage
	} else {
		response.ErrCode = errors.None
		response.ErrMessage = nil
		response.NodeID = kafkaResponse.NodeID
		response.Host = broker.AdvertiseListener.IP.String()
		response.Port = int32(broker.AdvertiseListener.Port)
	}

	return response, nil
}
