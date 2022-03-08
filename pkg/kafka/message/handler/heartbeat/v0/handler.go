package v0

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v0"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("heartbeat-v0-handler")

	request := message.(*v0.Request)

	response := &v0.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaHeartbeatRequest := &sarama.HeartbeatRequest{
		GroupId:      request.GroupID,
		GenerationId: request.GenerationID,
		MemberId:     request.MemberID,
	}

	kafkaHeartbeatResponse, err := client.Broker.GetController().HeartBeat(kafkaHeartbeatRequest)
	if err != nil {
		log.Error(err, "Error heartbeat to backend cluster")
		if kafkaError, ok := err.(sarama.KError); ok {
			response.ErrCode = errors.KafkaError(kafkaError)
		} else {
			return fmt.Errorf("error heartbeat to controller: %w", err)
		}
	}

	if kafkaHeartbeatResponse != nil {
		if response.ErrCode == errors.None {
			response.ErrCode = errors.KafkaError(kafkaHeartbeatResponse.Err)
		}
	}

	return client.WriteMessage(response, correlationId)
}
