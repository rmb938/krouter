package v1

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v1"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("list-offsets-v5-handler")
	request := message.(*v1.Request)

	response := &v1.Response{}

	for _, requestTopic := range request.Topics {
		topicResponse := v1.ListOffsetsTopicResponse{
			Name: requestTopic.Name,
		}

		cluster := broker.GetClusterByTopic(requestTopic.Name)

		for _, partition := range requestTopic.Partitions {
			partitionResponse := v1.ListOffsetsPartitionResponse{
				PartitionIndex: partition.PartitionIndex,
			}

			if cluster == nil {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			kafkaOffsetRequest := kmsg.NewPtrListOffsetsRequest()

			kafkaOffsetRequestTopic := kmsg.NewListOffsetsRequestTopic()
			kafkaOffsetRequestTopic.Topic = requestTopic.Name

			kafkaOffsetRequestTopicPartition := kmsg.NewListOffsetsRequestTopicPartition()
			kafkaOffsetRequestTopicPartition.Partition = partition.PartitionIndex
			kafkaOffsetRequestTopicPartition.Timestamp = partition.Timestamp.UnixMilli()

			kafkaOffsetRequestTopic.Partitions = append(kafkaOffsetRequestTopic.Partitions, kafkaOffsetRequestTopicPartition)

			kafkaOffsetRequest.Topics = append(kafkaOffsetRequest.Topics, kafkaOffsetRequestTopic)

			kafkaOffsetResponse, err := cluster.ListOffsets(broker.BrokerID, kafkaOffsetRequest)
			if err != nil {
				log.Error(err, "error listing offsets to backend cluster")
				return nil, fmt.Errorf("error listing offsets to backend cluster: %w", err)
			}

			if kafkaOffsetResponse != nil {
				block := kafkaOffsetResponse.Topics[0].Partitions[0]
				partitionResponse.ErrCode = errors.KafkaError(block.ErrorCode)
				partitionResponse.Timestamp = time.UnixMilli(block.Timestamp)
				partitionResponse.Offset = block.Offset
			}

			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Topics = append(response.Topics, topicResponse)
	}

	return response, nil
}
