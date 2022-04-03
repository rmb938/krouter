package v5

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v5 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v5"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("list-offsets-v5-handler")
	request := message.(*v5.Request)

	response := &v5.Response{}

	for _, requestTopic := range request.Topics {
		topicResponse := v5.ListOffsetsTopicResponse{
			Name: requestTopic.Name,
		}

		cluster := broker.GetClusterByTopic(requestTopic.Name)

		for _, partition := range requestTopic.Partitions {
			partitionResponse := v5.ListOffsetsPartitionResponse{
				PartitionIndex: partition.PartitionIndex,
			}

			if cluster == nil {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			kafkaOffsetRequest := kmsg.NewPtrListOffsetsRequest()
			kafkaOffsetRequest.IsolationLevel = request.IsolationLevel

			kafkaOffsetRequestTopic := kmsg.NewListOffsetsRequestTopic()
			kafkaOffsetRequestTopic.Topic = requestTopic.Name

			kafkaOffsetRequestTopicPartition := kmsg.NewListOffsetsRequestTopicPartition()
			kafkaOffsetRequestTopicPartition.Partition = partition.PartitionIndex
			kafkaOffsetRequestTopicPartition.CurrentLeaderEpoch = partition.CurrentLeaderEpoch
			kafkaOffsetRequestTopicPartition.Timestamp = partition.Timestamp.UnixMilli()

			kafkaOffsetRequestTopic.Partitions = append(kafkaOffsetRequestTopic.Partitions, kafkaOffsetRequestTopicPartition)

			kafkaOffsetRequest.Topics = append(kafkaOffsetRequest.Topics, kafkaOffsetRequestTopic)

			kafkaOffsetResponse, err := cluster.ListOffsets(broker.BrokerID, kafkaOffsetRequest)
			if err != nil {
				log.Error(err, "error listing offsets to backend cluster")
				return nil, fmt.Errorf("error listing offsets to backend cluster: %w", err)
			}

			if kafkaOffsetResponse != nil {
				if int64(kafkaOffsetResponse.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
					response.ThrottleDuration = time.Duration(kafkaOffsetResponse.ThrottleMillis) * time.Millisecond
				}

				block := kafkaOffsetResponse.Topics[0].Partitions[0]
				partitionResponse.ErrCode = errors.KafkaError(block.ErrorCode)
				partitionResponse.Timestamp = time.UnixMilli(block.Timestamp)
				partitionResponse.Offset = block.Offset
				partitionResponse.LeaderEpoch = block.LeaderEpoch
			}

			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Topics = append(response.Topics, topicResponse)
	}

	return response, nil
}
