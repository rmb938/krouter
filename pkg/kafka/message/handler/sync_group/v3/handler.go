package v3

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/franz"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("sync-group-v3-handler")

	request := message.(*v3.Request)

	response := &v3.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaSyncGroupRequest := franz.NewPtrSyncGroupRequest()
	kafkaSyncGroupRequest.Group = request.GroupID
	kafkaSyncGroupRequest.Generation = request.GenerationID
	kafkaSyncGroupRequest.MemberID = request.MemberID
	kafkaSyncGroupRequest.InstanceID = request.GroupInstanceID

	for _, assignment := range request.Assignments {
		kafkaSyncGroupRequestAssignment := franz.NewSyncGroupRequestGroupAssignment()
		kafkaSyncGroupRequestAssignment.MemberID = assignment.MemberID
		kafkaSyncGroupRequestAssignment.MemberAssignment = assignment.Assignment

		kafkaSyncGroupRequest.GroupAssignment = append(kafkaSyncGroupRequest.GroupAssignment, kafkaSyncGroupRequestAssignment)
	}

	kafkaSyncGroupResponse, err := broker.GetController().SyncGroupCustom(kafkaSyncGroupRequest)
	if err != nil {
		log.Error(err, "Error syncing group to backend cluster")
		return nil, fmt.Errorf("error syncing group to controller: %w", err)
	}

	response.ThrottleDuration = time.Duration(kafkaSyncGroupResponse.ThrottleMillis) * time.Millisecond
	response.ErrCode = errors.KafkaError(kafkaSyncGroupResponse.ErrorCode)
	response.Assignment = kafkaSyncGroupResponse.MemberAssignment

	return response, nil
}
