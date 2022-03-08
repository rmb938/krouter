package v3

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v3"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("list-offsets-v3-handler")
	request := message.(*v3.Request)

	response := &v3.Response{}

	for _, requestTopic := range request.Topics {
		topicResponse := v3.ListOffsetsTopicResponse{
			Name: requestTopic.Name,
		}

		cluster, topic := client.Broker.GetTopic(requestTopic.Name)

		for _, partition := range requestTopic.Partitions {
			partitionResponse := v3.ListOffsetsPartitionResponse{
				PartitionIndex: partition.PartitionIndex,
			}

			if topic == nil {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			if partition.PartitionIndex >= topic.Partitions {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			kafkaOffsetRequest := &sarama.OffsetRequest{
				Version:        request.Version(),
				IsolationLevel: sarama.IsolationLevel(request.IsolationLevel),
			}

			kafkaOffsetRequest.SetReplicaID(request.ReplicaID)
			kafkaOffsetRequest.AddBlock(requestTopic.Name, partition.PartitionIndex, partition.Timestamp.UnixMilli(), -1)

			kafkaOffsetResponse, err := cluster.ListOffsets(topic, partition.PartitionIndex, kafkaOffsetRequest)
			if err != nil {
				log.Error(err, "error listing offsets to backend cluster")
				if kafkaError, ok := err.(sarama.KError); ok {
					partitionResponse.ErrCode = errors.KafkaError(kafkaError)
				} else {
					return fmt.Errorf("error listing offsets to backend cluster: %w", err)
				}
			}

			if kafkaOffsetResponse != nil {
				if time.Duration(kafkaOffsetResponse.ThrottleTimeMs)*time.Millisecond > response.ThrottleDuration {
					response.ThrottleDuration = time.Duration(kafkaOffsetResponse.ThrottleTimeMs) * time.Millisecond
				}

				block := kafkaOffsetResponse.GetBlock(requestTopic.Name, partition.PartitionIndex)
				partitionResponse.Offset = block.Offset
				partitionResponse.Timestamp = time.UnixMilli(block.Timestamp)
			}

			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Topics = append(response.Topics, topicResponse)
	}

	return client.WriteMessage(response, correlationId)
}
