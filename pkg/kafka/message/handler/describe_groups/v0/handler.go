package v0

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("describe-groups-v0-handler")

	request := message.(*v0.Request)

	response := &v0.Response{}

	for _, group := range request.Groups {
		log = log.WithValues("group", group)

		responseGroup := v0.DescribeGroupGroupResponse{
			ErrCode: errors.GroupIDNotFound,
			GroupID: group,
		}

		kafkaResponseDescribeGroup, err := client.Broker.GetController().DescribeGroup(group)
		if err != nil {
			log.Error(err, "Error describing group to controller")
			if kafkaError, ok := err.(sarama.KError); ok {
				responseGroup.ErrCode = errors.KafkaError(kafkaError)
				response.Groups = append(response.Groups, responseGroup)
				continue
			}
			return fmt.Errorf("error describing group to controller: %w", err)
		}

		if kafkaResponseDescribeGroup != nil {
			kafkaResponseGroup := kafkaResponseDescribeGroup.Groups[0]
			responseGroup.ErrCode = errors.KafkaError(kafkaResponseGroup.Err)
			responseGroup.GroupState = kafkaResponseGroup.State
			responseGroup.ProtocolType = kafkaResponseGroup.ProtocolType
			responseGroup.ProtocolData = kafkaResponseGroup.Protocol

			for memberId, kafkaResponseMember := range kafkaResponseGroup.Members {
				responseGroup.Members = append(responseGroup.Members, v0.DescribeGroupMembersResponse{
					MemberID:         memberId,
					ClientID:         kafkaResponseMember.ClientId,
					ClientHost:       kafkaResponseMember.ClientHost,
					MemberMetadata:   kafkaResponseMember.MemberMetadata,
					MemberAssignment: kafkaResponseMember.MemberAssignment,
				})
			}

		}

		response.Groups = append(response.Groups, responseGroup)
	}

	return client.WriteMessage(response, correlationId)
}
