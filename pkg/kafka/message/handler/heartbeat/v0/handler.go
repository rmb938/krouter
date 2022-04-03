package v0

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v0"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("heartbeat-v0-handler")

	request := message.(*v0.Request)

	response := &v0.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaHeartbeatRequest := kmsg.NewPtrHeartbeatRequest()
	kafkaHeartbeatRequest.Group = request.GroupID
	kafkaHeartbeatRequest.Generation = request.GenerationID
	kafkaHeartbeatRequest.MemberID = request.MemberID

	kafkaHeartbeatResponse, err := broker.GetController().HeartBeat(kafkaHeartbeatRequest)
	if err != nil {
		log.Error(err, "Error heartbeat to backend cluster")
		return nil, fmt.Errorf("error heartbeat to controller: %w", err)
	}

	if kafkaHeartbeatResponse != nil {
		if response.ErrCode == errors.None {
			response.ErrCode = errors.KafkaError(kafkaHeartbeatResponse.ErrorCode)
		}
	}

	return response, nil
}
