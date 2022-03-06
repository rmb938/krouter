package v0

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v0"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("leave-group-v0-handler")
	request := message.(*v0.Request)

	response := &v0.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaLeaveGroupRequest := &sarama.LeaveGroupRequest{
		GroupId:  request.GroupID,
		MemberId: request.MemberID,
	}

	kafkaLeaveGroupResponse, err := client.Broker.GetController().LeaveGroup(kafkaLeaveGroupRequest)
	if err != nil {
		log.Error(err, "Error leaving group to backend cluster")
		if kafkaError, ok := err.(sarama.KError); ok {
			response.ErrCode = errors.KafkaError(kafkaError)
		} else {
			return fmt.Errorf("error leaving group to backend cluster: %w", err)
		}
	}

	if kafkaLeaveGroupResponse != nil {
		response.ErrCode = errors.KafkaError(kafkaLeaveGroupResponse.Err)
	}

	return client.WriteMessage(response, correlationId)
}
