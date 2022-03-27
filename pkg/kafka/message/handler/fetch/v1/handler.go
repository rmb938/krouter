package v1

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v1"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("fetch-v1-handler")

	request := message.(*v1.Request)

	response := &v1.Response{}

	for _, requestedTopic := range request.Topics {
		log = log.WithValues("topic", requestedTopic.Name)

		topicResponse := v1.FetchTopicResponse{
			Topic: requestedTopic.Name,
		}

		cluster := broker.GetClusterByTopic(requestedTopic.Name)

		for _, partition := range requestedTopic.Partitions {
			log = log.WithValues("partition", partition.Partition)

			partitionResponse := v1.FetchPartitionResponse{
				PartitionIndex: partition.Partition,
			}

			if cluster == nil {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			kafkaFetchRequest := kmsg.NewPtrFetchRequest()
			kafkaFetchRequest.ReplicaID = -1
			kafkaFetchRequest.MaxWaitMillis = int32(request.MaxWait.Milliseconds())
			kafkaFetchRequest.MinBytes = request.MinBytes

			kafkaFetchRequestTopic := kmsg.NewFetchRequestTopic()
			kafkaFetchRequestTopic.Topic = requestedTopic.Name

			kafkaFetchRequestTopicPartition := kmsg.NewFetchRequestTopicPartition()
			kafkaFetchRequestTopicPartition.Partition = partition.Partition
			kafkaFetchRequestTopicPartition.FetchOffset = partition.FetchOffset
			kafkaFetchRequestTopicPartition.PartitionMaxBytes = partition.PartitionMaxBytes

			kafkaFetchRequestTopic.Partitions = append(kafkaFetchRequestTopic.Partitions, kafkaFetchRequestTopicPartition)

			kafkaFetchRequest.Topics = append(kafkaFetchRequest.Topics, kafkaFetchRequestTopic)

			kafkaFetchResponse, err := cluster.Fetch(requestedTopic.Name, partition.Partition, kafkaFetchRequest)
			if err != nil {
				log.Error(err, "Error fetch to backend cluster")
				return nil, fmt.Errorf("error fetchto backend cluster: %w", err)
			}

			if int64(kafkaFetchResponse.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
				response.ThrottleDuration = time.Duration(kafkaFetchResponse.ThrottleMillis) * time.Millisecond
			}

			block := kafkaFetchResponse.Topics[0].Partitions[0]

			partitionResponse.ErrCode = errors.KafkaError(block.ErrorCode)
			partitionResponse.HighWaterMark = block.HighWatermark

			responseRecordByteBuff := bytes.NewBuffer([]byte{})
			responseRecordByteBuff.Write(block.RecordBatches)

			responseRecordByteBuffLen := responseRecordByteBuff.Len()
			if responseRecordByteBuffLen > 0 {
				partitionResponse.Records = responseRecordByteBuff.Bytes()
			}

			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Responses = append(response.Responses, topicResponse)
	}

	return response, nil
}
