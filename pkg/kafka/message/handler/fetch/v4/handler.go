package v4

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v4"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("fetch-v4-handler")

	request := message.(*v4.Request)

	response := &v4.Response{}

	responseSize := int64(0)

	clusterRequests := make(map[*logical_broker.Cluster][]kmsg.FetchRequestTopic)
	for _, requestedTopic := range request.Topics {
		cluster := broker.GetClusterByTopic(requestedTopic.Name)

		if cluster == nil {
			topicResponse := v4.FetchTopicResponse{
				Topic: requestedTopic.Name,
			}

			for _, partition := range requestedTopic.Partitions {
				partitionResponse := v4.FetchPartitionResponse{
					PartitionIndex: partition.Partition,
					ErrCode:        errors.UnknownTopicOrPartition,
				}
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
			}

			response.Responses = append(response.Responses, topicResponse)
			continue
		}

		topicRequests, ok := clusterRequests[cluster]
		if !ok {
			topicRequests = make([]kmsg.FetchRequestTopic, 0)
		}

		kafkaFetchRequestTopic := kmsg.NewFetchRequestTopic()
		kafkaFetchRequestTopic.Topic = requestedTopic.Name

		for _, partition := range requestedTopic.Partitions {
			kafkaFetchRequestTopicPartition := kmsg.NewFetchRequestTopicPartition()
			kafkaFetchRequestTopicPartition.Partition = partition.Partition
			kafkaFetchRequestTopicPartition.FetchOffset = partition.FetchOffset
			kafkaFetchRequestTopicPartition.PartitionMaxBytes = partition.PartitionMaxBytes
			kafkaFetchRequestTopic.Partitions = append(kafkaFetchRequestTopic.Partitions, kafkaFetchRequestTopicPartition)
		}

		topicRequests = append(topicRequests, kafkaFetchRequestTopic)

		clusterRequests[cluster] = topicRequests
	}

	for cluster, topicRequests := range clusterRequests {
		if responseSize >= int64(request.MaxBytes) {
			continue
		}

		kafkaFetchRequest := kmsg.NewPtrFetchRequest()
		kafkaFetchRequest.ReplicaID = -1
		kafkaFetchRequest.MaxWaitMillis = int32(request.MaxWait.Milliseconds())
		kafkaFetchRequest.MinBytes = request.MinBytes
		kafkaFetchRequest.MaxBytes = request.MaxBytes
		kafkaFetchRequest.IsolationLevel = request.IsolationLevel
		kafkaFetchRequest.Topics = topicRequests

		kafkaFetchResponse, err := cluster.Fetch(broker.BrokerID, kafkaFetchRequest)
		if err != nil {
			log.Error(err, "Error fetch to backend cluster")
			return nil, fmt.Errorf("error fetchto backend cluster: %w", err)
		}

		if int64(kafkaFetchResponse.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
			response.ThrottleDuration = time.Duration(kafkaFetchResponse.ThrottleMillis) * time.Millisecond
		}

		for _, responseTopic := range kafkaFetchResponse.Topics {
			topicResponse := v4.FetchTopicResponse{
				Topic: responseTopic.Topic,
			}

			for _, responsePartition := range responseTopic.Partitions {
				partitionResponse := v4.FetchPartitionResponse{
					PartitionIndex:   responsePartition.Partition,
					ErrCode:          errors.KafkaError(responsePartition.ErrorCode),
					HighWaterMark:    responsePartition.HighWatermark,
					LastStableOffset: responsePartition.LastStableOffset,
					Records:          responsePartition.RecordBatches,
				}

				responseSize += int64(len(partitionResponse.Records))

				for _, responseAbortedTransaction := range responsePartition.AbortedTransactions {
					partitionResponse.AbortedTransactions = append(partitionResponse.AbortedTransactions, v4.FetchAbortedTransaction{
						ProducerID:  responseAbortedTransaction.ProducerID,
						FirstOffset: responseAbortedTransaction.FirstOffset,
					})
				}

				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
			}

			response.Responses = append(response.Responses, topicResponse)
		}
	}

	return response, nil
}
