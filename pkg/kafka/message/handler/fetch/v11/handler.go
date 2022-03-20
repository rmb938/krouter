package v11

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("fetch-v11-handler")

	request := message.(*v11.Request)

	response := &v11.Response{
		SessionID: request.SessionID,
	}

	responseSize := int64(0)

topicLoop:
	for _, requestedTopic := range request.Topics {
		log = log.WithValues("topic", requestedTopic.Name)

		if responseSize >= int64(request.MaxBytes) {
			continue
		}

		topicResponse := v11.FetchTopicResponse{
			Topic: requestedTopic.Name,
		}

		cluster, topic := broker.GetTopic(requestedTopic.Name)

		for _, partition := range requestedTopic.Partitions {
			log = log.WithValues("partition", partition.Partition)

			if responseSize >= int64(request.MaxBytes) {
				continue
			}

			partitionResponse := v11.FetchPartitionResponse{
				PartitionIndex:       partition.Partition,
				PreferredReadReplica: 1,
			}

			if topic == nil {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			kafkaFetchRequest := kmsg.NewPtrFetchRequest()
			kafkaFetchRequest.ReplicaID = -1
			kafkaFetchRequest.MaxWaitMillis = int32(request.MaxWait.Milliseconds())
			kafkaFetchRequest.MinBytes = request.MinBytes
			kafkaFetchRequest.MaxBytes = request.MaxBytes
			kafkaFetchRequest.IsolationLevel = request.IsolationLevel
			kafkaFetchRequest.SessionID = request.SessionID
			kafkaFetchRequest.SessionEpoch = request.SessionEpoch
			kafkaFetchRequest.Rack = request.RackID

			kafkaFetchRequestTopic := kmsg.NewFetchRequestTopic()
			kafkaFetchRequestTopic.Topic = topic.Name

			kafkaFetchRequestTopicPartition := kmsg.NewFetchRequestTopicPartition()
			kafkaFetchRequestTopicPartition.Partition = partition.Partition
			kafkaFetchRequestTopicPartition.CurrentLeaderEpoch = partition.CurrentLeaderEpoch
			kafkaFetchRequestTopicPartition.FetchOffset = partition.FetchOffset
			kafkaFetchRequestTopicPartition.LogStartOffset = partition.LogStartOffset
			kafkaFetchRequestTopicPartition.PartitionMaxBytes = partition.PartitionMaxBytes

			kafkaFetchRequestTopic.Partitions = append(kafkaFetchRequestTopic.Partitions, kafkaFetchRequestTopicPartition)

			kafkaFetchRequest.Topics = append(kafkaFetchRequest.Topics, kafkaFetchRequestTopic)

			kafkaFetchResponse, err := cluster.Fetch(topic, partition.Partition, kafkaFetchRequest)
			if err != nil {
				log.Error(err, "Error fetch to backend cluster")
				return nil, fmt.Errorf("error fetchto backend cluster: %w", err)
			}

			if int64(kafkaFetchResponse.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
				response.ThrottleDuration = time.Duration(kafkaFetchResponse.ThrottleMillis) * time.Millisecond
			}

			// something returned a high level error so stop everything and just return nothing
			if kafkaFetchResponse.ErrorCode != int16(errors.None) {
				response.ErrCode = errors.KafkaError(kafkaFetchResponse.ErrorCode)
				response.Responses = nil
				break topicLoop
			}

			block := kafkaFetchResponse.Topics[0].Partitions[0]

			partitionResponse.ErrCode = errors.KafkaError(block.ErrorCode)
			partitionResponse.HighWaterMark = block.HighWatermark
			partitionResponse.LastStableOffset = block.LastStableOffset
			partitionResponse.LogStartOffset = block.LogStartOffset

			for _, blockAbortedTransaction := range block.AbortedTransactions {
				partitionResponse.AbortedTransactions = append(partitionResponse.AbortedTransactions, v11.FetchAbortedTransaction{
					ProducerID:  blockAbortedTransaction.ProducerID,
					FirstOffset: blockAbortedTransaction.FirstOffset,
				})
			}

			responseRecordByteBuff := bytes.NewBuffer([]byte{})
			responseRecordByteBuff.Write(block.RecordBatches)

			responseRecordByteBuffLen := responseRecordByteBuff.Len()
			if responseRecordByteBuffLen > 0 {
				responseSize += int64(responseRecordByteBuffLen)
				partitionResponse.Records = responseRecordByteBuff.Bytes()
			}

			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Responses = append(response.Responses, topicResponse)
	}

	return response, nil
}
