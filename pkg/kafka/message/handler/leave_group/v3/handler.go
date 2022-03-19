package v3

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v3"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("leave-group-v0-handler")
	request := message.(*v3.Request)

	response := &v3.Response{}

	log = log.WithValues("group-id", request.GroupID)

	kafkaLeaveGroupRequest := kmsg.NewPtrLeaveGroupRequest()
	kafkaLeaveGroupRequest.Group = request.GroupID
	for _, member := range request.Members {
		kafkaLeaveGroupRequestMember := kmsg.NewLeaveGroupRequestMember()
		kafkaLeaveGroupRequestMember.MemberID = member.MemberID
		kafkaLeaveGroupRequestMember.InstanceID = member.GroupInstanceID

		kafkaLeaveGroupRequest.Members = append(kafkaLeaveGroupRequest.Members, kafkaLeaveGroupRequestMember)
	}

	kafkaLeaveGroupResponse, err := client.Broker.GetController().LeaveGroup(kafkaLeaveGroupRequest)
	if err != nil {
		log.Error(err, "Error leaving group to controller")
		return fmt.Errorf("error leaving group to backend cluster: %w", err)
	}

	response.ThrottleDuration = time.Duration(kafkaLeaveGroupResponse.ThrottleMillis) * time.Millisecond
	response.ErrCode = errors.KafkaError(kafkaLeaveGroupResponse.ErrorCode)

	for _, member := range kafkaLeaveGroupResponse.Members {
		response.Members = append(response.Members, v3.ResponseMember{
			MemberID:        member.MemberID,
			GroupInstanceID: member.InstanceID,
			ErrCode:         errors.KafkaError(member.ErrorCode),
		})
	}

	return client.WriteMessage(response, correlationId)
}
