package v0

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v5 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v5"
	"github.com/twmb/franz-go/pkg/kmsg"

	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("join-group-v4-handler")

	request := message.(*v5.Request)

	response := &v5.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	kafkaJoinGroupRequest := kmsg.NewPtrJoinGroupRequest()
	kafkaJoinGroupRequest.Group = request.GroupID
	kafkaJoinGroupRequest.SessionTimeoutMillis = int32(request.SessionTimeout.Milliseconds())
	kafkaJoinGroupRequest.RebalanceTimeoutMillis = int32(request.RebalanceTimeout.Milliseconds())
	kafkaJoinGroupRequest.MemberID = request.MemberID
	kafkaJoinGroupRequest.InstanceID = request.GroupInstanceId
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

	if int64(kafkaJoinGroupResponse.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
		response.ThrottleDuration = time.Duration(kafkaJoinGroupResponse.ThrottleMillis) * time.Millisecond
	}

	response.ErrCode = errors.KafkaError(kafkaJoinGroupResponse.ErrorCode)

	response.GenerationID = kafkaJoinGroupResponse.Generation

	if kafkaJoinGroupResponse.Protocol != nil {
		response.ProtocolName = *kafkaJoinGroupResponse.Protocol
	}

	response.Leader = kafkaJoinGroupResponse.LeaderID
	response.MemberID = kafkaJoinGroupResponse.MemberID

	for _, kafkaMemberMetadata := range kafkaJoinGroupResponse.Members {
		response.Members = append(response.Members, v5.GroupMember{
			MemberID:        kafkaMemberMetadata.MemberID,
			GroupInstanceID: kafkaMemberMetadata.InstanceID,
			Metadata:        kafkaMemberMetadata.ProtocolMetadata,
		})
	}

	return response, nil
}
