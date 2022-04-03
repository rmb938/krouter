package v0

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v0"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("leave-group-v0-handler")
	request := message.(*v0.Request)

	response := &v0.Response{}

	log = log.WithValues("group-id", request.GroupID)

	kafkaLeaveGroupRequest := kmsg.NewPtrLeaveGroupRequest()
	kafkaLeaveGroupRequest.Group = request.GroupID

	// Franz always uses the latest version, so we need to use the latest as well
	kafkaLeaveGroupRequestMember := kmsg.NewLeaveGroupRequestMember()
	kafkaLeaveGroupRequestMember.MemberID = request.MemberID
	kafkaLeaveGroupRequest.Members = append(kafkaLeaveGroupRequest.Members, kafkaLeaveGroupRequestMember)

	kafkaLeaveGroupResponse, err := broker.GetController().LeaveGroup(kafkaLeaveGroupRequest)
	if err != nil {
		log.Error(err, "Error leaving group to controller")
		return nil, fmt.Errorf("error leaving group to backend cluster: %w", err)
	}

	response.ErrCode = errors.KafkaError(kafkaLeaveGroupResponse.ErrorCode)

	return response, nil
}
