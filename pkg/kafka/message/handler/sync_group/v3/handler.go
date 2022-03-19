package v3

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("sync-group-v0-handler")

	request := message.(*v3.Request)

	response := &v3.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaSyncGroupRequest := kmsg.NewPtrSyncGroupRequest()
	kafkaSyncGroupRequest.Group = request.GroupID
	kafkaSyncGroupRequest.Generation = request.GenerationID
	kafkaSyncGroupRequest.MemberID = request.MemberID
	kafkaSyncGroupRequest.InstanceID = request.GroupInstanceID

	for _, assignment := range request.Assignments {
		kafkaSyncGroupRequestAssignment := kmsg.NewSyncGroupRequestGroupAssignment()
		kafkaSyncGroupRequestAssignment.MemberID = assignment.MemberID
		kafkaSyncGroupRequestAssignment.MemberAssignment = assignment.Assignment

		kafkaSyncGroupRequest.GroupAssignment = append(kafkaSyncGroupRequest.GroupAssignment, kafkaSyncGroupRequestAssignment)
	}

	kafkaSyncGroupResponse, err := client.Broker.GetController().SyncGroup(kafkaSyncGroupRequest)
	if err != nil {
		log.Error(err, "Error syncing group to backend cluster")
		return fmt.Errorf("error syncing group to controller: %w", err)
	}

	response.ThrottleDuration = time.Duration(kafkaSyncGroupResponse.ThrottleMillis) * time.Millisecond
	response.ErrCode = errors.KafkaError(kafkaSyncGroupResponse.ErrorCode)
	response.Assignment = kafkaSyncGroupResponse.MemberAssignment

	return client.WriteMessage(response, correlationId)
}
