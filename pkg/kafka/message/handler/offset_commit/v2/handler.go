package v2

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v2"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("offset-commit-v2-handler")
	request := message.(*v2.Request)

	response := &v2.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	for _, requestTopic := range request.Topics {
		topicResponse := v2.OffsetCommitTopicResponse{
			Name: requestTopic.Name,
		}

		_, topic := broker.GetTopic(requestTopic.Name)

		for _, requestPartition := range requestTopic.Partitions {
			partitionResponse := v2.OffsetCommitPartitionResponse{
				PartitionIndex: requestPartition.PartitionIndex,
				ErrCode:        errors.None,
			}

			if topic == nil {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			if requestPartition.PartitionIndex >= topic.Partitions {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			kafkaErr, err := broker.GetController().OffsetCommit(request.GroupID, topic.Name, request.GenerationID, requestPartition.PartitionIndex, requestPartition.CommittedOffset)
			if err != nil {
				log.Error(err, "Error offset commit to controller")
				return nil, fmt.Errorf("error offset commit to backend cluster: %w", err)
			}

			partitionResponse.ErrCode = kafkaErr
			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Topics = append(response.Topics, topicResponse)
	}

	return response, nil
}
