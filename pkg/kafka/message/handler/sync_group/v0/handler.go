package v0

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v0"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("sync-group-v0-handler")

	request := message.(*v0.Request)

	response := &v0.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaSyncGroupRequest := &sarama.SyncGroupRequest{
		GroupId:          request.GroupID,
		GenerationId:     request.GenerationID,
		MemberId:         request.MemberID,
		GroupAssignments: map[string][]byte{},
	}

	for _, groupAssignment := range request.Assignments {
		kafkaSyncGroupRequest.GroupAssignments[groupAssignment.MemberID] = groupAssignment.Assignment
	}

	kafkaSyncGroupResponse, err := client.Broker.GetController().SyncGroup(kafkaSyncGroupRequest)
	if err != nil {
		log.Error(err, "Error syncing group to backend cluster")
		if kafkaError, ok := err.(sarama.KError); ok {
			response.ErrCode = errors.KafkaError(kafkaError)
		} else {
			return fmt.Errorf("error syncing group to controller: %w", err)
		}
	}

	if kafkaSyncGroupResponse != nil {
		response.ErrCode = errors.KafkaError(kafkaSyncGroupResponse.Err)
		response.Assignment = kafkaSyncGroupResponse.MemberAssignment
	}

	return client.WriteMessage(response, correlationId)
}
