package v0

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v0"
	"github.com/twmb/franz-go/pkg/kmsg"

	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("join-group-v0-handler")

	request := message.(*v0.Request)

	response := &v0.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaJoinGroupRequest := kmsg.NewPtrJoinGroupRequest()
	kafkaJoinGroupRequest.Group = request.GroupID
	kafkaJoinGroupRequest.SessionTimeoutMillis = int32(request.SessionTimeout.Milliseconds())
	kafkaJoinGroupRequest.MemberID = request.MemberID
	kafkaJoinGroupRequest.ProtocolType = request.ProtocolType

	for _, protocol := range request.Protocols {
		kafkaJoinGroupRequestProtocol := kmsg.NewJoinGroupRequestProtocol()
		kafkaJoinGroupRequestProtocol.Name = protocol.Name
		kafkaJoinGroupRequestProtocol.Metadata = protocol.Metadata

		kafkaJoinGroupRequest.Protocols = append(kafkaJoinGroupRequest.Protocols, kafkaJoinGroupRequestProtocol)
	}

	kafkaJoinGroupResponse, err := broker.GetController().JoinGroup(kafkaJoinGroupRequest)
	if err != nil {
		log.Error(err, "Error joining group to backend cluster")
		return nil, fmt.Errorf("error joining group to controller: %w", err)
	}

	response.ErrCode = errors.KafkaError(kafkaJoinGroupResponse.ErrorCode)

	response.GenerationID = kafkaJoinGroupResponse.Generation

	if kafkaJoinGroupResponse.Protocol != nil {
		response.ProtocolName = *kafkaJoinGroupResponse.Protocol
	}

	response.Leader = kafkaJoinGroupResponse.LeaderID
	response.MemberID = kafkaJoinGroupResponse.MemberID

	for _, kafkaMemberMetadata := range kafkaJoinGroupResponse.Members {
		response.Members = append(response.Members, v0.GroupMember{
			MemberID: kafkaMemberMetadata.MemberID,
			Metadata: kafkaMemberMetadata.ProtocolMetadata,
		})
	}

	return response, nil
}
