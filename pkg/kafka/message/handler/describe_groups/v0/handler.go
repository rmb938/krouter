package v0

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("describe-groups-v0-handler")

	request := message.(*v0.Request)

	response := &v0.Response{}

	for _, group := range request.Groups {
		log = log.WithValues("group", group)

		responseGroup := v0.DescribeGroupGroupResponse{
			ErrCode: errors.GroupIDNotFound,
			GroupID: group,
		}

		kafkaResponseDescribeGroup, err := broker.GetController().DescribeGroup(group)
		if err != nil {
			log.Error(err, "Error describing group to controller")
			return nil, fmt.Errorf("error describing group to controller: %w", err)
		}

		kafkaResponseGroup := kafkaResponseDescribeGroup.Groups[0]
		responseGroup.ErrCode = errors.KafkaError(kafkaResponseGroup.ErrorCode)
		responseGroup.GroupState = kafkaResponseGroup.State
		responseGroup.ProtocolType = kafkaResponseGroup.ProtocolType
		responseGroup.ProtocolData = kafkaResponseGroup.Protocol

		for _, kafkaResponseMember := range kafkaResponseGroup.Members {
			responseGroup.Members = append(responseGroup.Members, v0.DescribeGroupMembersResponse{
				MemberID:         kafkaResponseMember.MemberID,
				ClientID:         kafkaResponseMember.ClientID,
				ClientHost:       kafkaResponseMember.ClientHost,
				MemberMetadata:   kafkaResponseMember.ProtocolMetadata,
				MemberAssignment: kafkaResponseMember.MemberAssignment,
			})
		}

		response.Groups = append(response.Groups, responseGroup)
	}

	return response, nil
}
