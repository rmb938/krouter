package v5

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v5 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v5"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("offset-fetch-v5-handler")

	request := message.(*v5.Request)

	log = log.WithValues("group-id", request.GroupID)

	response := &v5.Response{}

	kafkaOffsetFetchRequest := &sarama.OffsetFetchRequest{
		Version:       request.Version(),
		ConsumerGroup: request.GroupID,
	}

	kafkaOffsetFetchRequest.ZeroPartitions()

	for _, topic := range request.Topics {
		for _, partitionIndex := range topic.PartitionIndexes {
			kafkaOffsetFetchRequest.AddPartition(topic.Name, partitionIndex)
		}
	}

	kafkaOffsetFetchResponse, err := client.Broker.GetController().OffsetFetch(kafkaOffsetFetchRequest)
	if err != nil {
		log.Error(err, "Error syncing group to backend cluster")
		if kafkaError, ok := err.(sarama.KError); ok {
			response.ErrCode = errors.KafkaError(kafkaError)
		} else {
			return fmt.Errorf("error syncing group to backend cluster: %w", err)
		}
	}

	if kafkaOffsetFetchResponse != nil {
		if time.Duration(kafkaOffsetFetchResponse.ThrottleTimeMs)*time.Millisecond > response.ThrottleDuration {
			response.ThrottleDuration = time.Duration(kafkaOffsetFetchResponse.ThrottleTimeMs) * time.Millisecond
		}

		for topic, partitions := range kafkaOffsetFetchResponse.Blocks {
			offsetFetchTopic := v5.ResponseOffsetFetchTopic{
				Name: topic,
			}

			for partitionIndex, partitionInfo := range partitions {
				var metadata *string
				if len(partitionInfo.Metadata) > 0 {
					metadata = &partitionInfo.Metadata
				}
				offsetFetchTopicPartition := v5.ResponseOffsetFetchTopicPartition{
					PartitionIndex:       partitionIndex,
					CommittedOffset:      partitionInfo.Offset,
					CommittedLeaderEpoch: partitionInfo.LeaderEpoch,
					Metadata:             metadata,
					ErrCode:              errors.KafkaError(partitionInfo.Err),
				}

				offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, offsetFetchTopicPartition)
			}

			response.Topics = append(response.Topics, offsetFetchTopic)
		}

		response.ErrCode = errors.KafkaError(kafkaOffsetFetchResponse.Err)
	}

	return client.WriteMessage(response, correlationId)
}
