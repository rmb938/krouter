package v4

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v4"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("offset-commit-v4-handler")
	request := message.(*v4.Request)

	response := &v4.Response{}

	log = log.WithValues("group-id", request.GroupID, "member-id", request.MemberID)

	for _, requestTopic := range request.Topics {
		topicResponse := v4.OffsetCommitTopicResponse{
			Name: requestTopic.Name,
		}

		_, topic := client.Broker.GetTopic(requestTopic.Name)

		for _, requestPartition := range requestTopic.Partitions {
			partitionResponse := v4.OffsetCommitPartitionResponse{
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

			err := client.Broker.GetController().OffsetCommit(request.GroupID, topic.Name, request.GenerationID, requestPartition.PartitionIndex, requestPartition.CommittedOffset)
			if err != nil {
				log.Error(err, "Error offset commit to controller")
				if kafkaError, ok := err.(sarama.KError); ok {
					partitionResponse.ErrCode = errors.KafkaError(kafkaError)
				} else {
					return fmt.Errorf("error offset commit to backend cluster: %w", err)
				}
			}

			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Topics = append(response.Topics, topicResponse)
	}

	return client.WriteMessage(response, correlationId)
}
