package v4

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("offset-fetch-v5-handler")

	request := message.(*v4.Request)

	log = log.WithValues("group-id", request.GroupID)

	response := &v4.Response{
		ErrCode: errors.None,
	}

	for _, requestTopic := range request.Topics {
		offsetFetchTopic := v4.ResponseOffsetFetchTopic{
			Name: requestTopic.Name,
		}

		_, topic := client.Broker.GetTopic(requestTopic.Name)

		for _, partitionIndex := range requestTopic.PartitionIndexes {
			offsetFetchTopicPartition := v4.ResponseOffsetFetchTopicPartition{
				PartitionIndex: partitionIndex,
				Metadata:       nil,
				ErrCode:        errors.None,
			}

			if topic == nil {
				offsetFetchTopicPartition.ErrCode = errors.UnknownTopicOrPartition
				offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, offsetFetchTopicPartition)
				continue
			}

			if partitionIndex >= topic.Partitions {
				offsetFetchTopicPartition.ErrCode = errors.UnknownTopicOrPartition
				offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, offsetFetchTopicPartition)
				continue
			}

			offset, err := client.Broker.GetController().OffsetFetch(request.GroupID, topic.Name, partitionIndex)
			if err != nil {
				log.Error(err, "error fetching offsets to controller")
				if kafkaError, ok := err.(sarama.KError); ok {
					offsetFetchTopicPartition.ErrCode = errors.KafkaError(kafkaError)
				} else {
					return fmt.Errorf("error fetching offsets to controller: %w", err)
				}
			}

			if offsetFetchTopicPartition.ErrCode == errors.None {
				offsetFetchTopicPartition.CommittedOffset = offset
			}

			offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, offsetFetchTopicPartition)
		}

		response.Topics = append(response.Topics, offsetFetchTopic)
	}

	return client.WriteMessage(response, correlationId)
}
