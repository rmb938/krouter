package v0

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v4"

	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("join-group-v4-handler")

	request := message.(*v4.Request)

	response := &v4.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaJoinGroupRequest := &sarama.JoinGroupRequest{
		Version:               request.Version(),
		GroupId:               request.GroupID,
		SessionTimeout:        int32(request.SessionTimeout.Milliseconds()),
		RebalanceTimeout:      int32(request.RebalanceTimeout.Milliseconds()),
		MemberId:              request.MemberID,
		ProtocolType:          request.ProtocolType,
		OrderedGroupProtocols: nil,
	}

	for _, protocol := range request.Protocols {
		kafkaGroupProtocol := &sarama.GroupProtocol{
			Name:     protocol.Name,
			Metadata: protocol.Metadata,
		}

		kafkaJoinGroupRequest.OrderedGroupProtocols = append(kafkaJoinGroupRequest.OrderedGroupProtocols, kafkaGroupProtocol)
	}

	kafkaJoinGroupResponse, err := client.Broker.GetController().JoinGroup(kafkaJoinGroupRequest)
	if err != nil {
		log.Error(err, "Error joining group to backend cluster")
		if kafkaError, ok := err.(sarama.KError); ok {
			response.ErrCode = errors.KafkaError(kafkaError)
		} else if err == sarama.ErrControllerNotAvailable {
			response.ErrCode = errors.UnknownServerError
		} else {
			return fmt.Errorf("error joining group to backend cluster: %w", err)
		}
	}

	if kafkaJoinGroupResponse != nil {
		if time.Duration(kafkaJoinGroupResponse.ThrottleTime) > response.ThrottleDuration {
			response.ThrottleDuration = time.Duration(kafkaJoinGroupResponse.ThrottleTime)
		}

		response.ErrCode = errors.KafkaError(kafkaJoinGroupResponse.Err)

		response.GenerationID = kafkaJoinGroupResponse.GenerationId
		response.ProtocolName = kafkaJoinGroupResponse.GroupProtocol
		response.Leader = kafkaJoinGroupResponse.LeaderId
		response.MemberID = kafkaJoinGroupResponse.MemberId

		for kafkaMemberID, kafkaMemberMetadata := range kafkaJoinGroupResponse.Members {
			response.Members = append(response.Members, v4.GroupMember{
				MemberID: kafkaMemberID,
				Metadata: kafkaMemberMetadata,
			})
		}
	}

	return client.WriteMessage(response, correlationId)
}
