package v3

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v3"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("heartbeat-v0-handler")

	request := message.(*v3.Request)

	response := &v3.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaHeartbeatRequest := kmsg.NewPtrHeartbeatRequest()
	kafkaHeartbeatRequest.Group = request.GroupID
	kafkaHeartbeatRequest.Generation = request.GenerationID
	kafkaHeartbeatRequest.MemberID = request.MemberID
	kafkaHeartbeatRequest.InstanceID = request.GroupInstanceId

	kafkaHeartbeatResponse, err := broker.GetController().HeartBeat(kafkaHeartbeatRequest)
	if err != nil {
		log.Error(err, "Error heartbeat to backend cluster")
		return nil, fmt.Errorf("error heartbeat to controller: %w", err)
	}

	response.ThrottleDuration = time.Duration(kafkaHeartbeatResponse.ThrottleMillis) * time.Millisecond

	if kafkaHeartbeatResponse != nil {
		if response.ErrCode == errors.None {
			response.ErrCode = errors.KafkaError(kafkaHeartbeatResponse.ErrorCode)
		}
	}

	return response, nil
}
