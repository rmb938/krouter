package v1

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("find-coordinator-v2-handler")

	request := message.(*v2.Request)

	response := &v2.Response{}

	if request.KeyType != 0 {
		response.ErrCode = errors.TransactionIDAuthorizationFailed
	} else {
		broker := client.Broker

		kafkaResponse, err := broker.GetController().FindCoordinator(request.Key)
		if err != nil {
			return err
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
	}

	return client.WriteMessage(response, correlationId)
}
