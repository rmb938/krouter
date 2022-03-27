package v3

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v3"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("produce-v7-handler")
	request := message.(*v3.Request)

	response := &v3.Response{}

	if len(request.TopicData) == 0 {
		return response, nil
	}

	response.Responses = make([]v3.ProduceResponse, 0, len(request.TopicData))

	if request.TransactionalID != nil {
		for _, topicData := range request.TopicData {
			produceResponse := v3.ProduceResponse{
				Name:               topicData.Name,
				PartitionResponses: make([]v3.PartitionResponse, 0, len(topicData.PartitionData)),
			}
			for _, partitionData := range topicData.PartitionData {
				partitionResponse := v3.PartitionResponse{
					Index:   partitionData.Index,
					ErrCode: errors.TransactionIDAuthorizationFailed,
				}
				produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
			}
			response.Responses = append(response.Responses, produceResponse)
		}

		return response, nil
	}

	if request.ACKs < -1 || request.ACKs > 1 {
		for _, topicData := range request.TopicData {
			produceResponse := v3.ProduceResponse{
				Name:               topicData.Name,
				PartitionResponses: make([]v3.PartitionResponse, 0, len(topicData.PartitionData)),
			}
			for _, partitionData := range topicData.PartitionData {
				partitionResponse := v3.PartitionResponse{
					Index:   partitionData.Index,
					ErrCode: errors.InvalidRequiredAcks,
				}
				produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
			}
			response.Responses = append(response.Responses, produceResponse)
		}

		return response, nil
	}

	topicRequestByCluster := make(map[*logical_broker.Cluster][]kmsg.ProduceRequestTopic)
	for _, topicData := range request.TopicData {
		cluster := broker.GetClusterByTopic(topicData.Name)

		if cluster == nil {
			produceResponse := v3.ProduceResponse{
				Name:               topicData.Name,
				PartitionResponses: make([]v3.PartitionResponse, 0, len(topicData.PartitionData)),
			}
			for _, partitionData := range topicData.PartitionData {
				partitionResponse := v3.PartitionResponse{
					Index:   partitionData.Index,
					ErrCode: errors.UnknownTopicOrPartition,
				}
				produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
			}
			response.Responses = append(response.Responses, produceResponse)
		}

		topicRequests, ok := topicRequestByCluster[cluster]
		if !ok {
			topicRequests = make([]kmsg.ProduceRequestTopic, 0)
		}

		topicRequest := kmsg.NewProduceRequestTopic()
		topicRequest.Topic = topicData.Name
		topicRequest.Partitions = make([]kmsg.ProduceRequestTopicPartition, 0, len(topicData.PartitionData))

		for _, partitionData := range topicData.PartitionData {
			produceTopicPartition := kmsg.NewProduceRequestTopicPartition()

			produceTopicPartition.Partition = partitionData.Index
			produceTopicPartition.Records = partitionData.Records

			topicRequest.Partitions = append(topicRequest.Partitions, produceTopicPartition)
		}

		topicRequestByCluster[cluster] = append(topicRequests, topicRequest)
	}

	for cluster, produceRequests := range topicRequestByCluster {
		kafkaResponse, err := cluster.Produce(request.TransactionalID, int32(request.TimeoutDuration.Milliseconds()), produceRequests)
		if err != nil {
			log.Error(err, "Error producing message to backend cluster")
			return nil, fmt.Errorf("error producing to kafka: %w", err)
		}

		if request.ACKs != 0 {
			if kafkaResponse.ThrottleMillis > int32(response.ThrottleDuration) {
				response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond
			}

			for _, kafkaRespTopic := range kafkaResponse.Topics {
				produceResponse := v3.ProduceResponse{
					Name:               kafkaRespTopic.Topic,
					PartitionResponses: make([]v3.PartitionResponse, 0, len(kafkaRespTopic.Partitions)),
				}

				for _, kafkaRespTopicPartition := range kafkaRespTopic.Partitions {
					produceResponsePartition := v3.PartitionResponse{
						Index:         kafkaRespTopicPartition.Partition,
						ErrCode:       errors.KafkaError(kafkaRespTopicPartition.ErrorCode),
						BaseOffset:    kafkaRespTopicPartition.BaseOffset,
						LogAppendTime: time.UnixMilli(kafkaRespTopicPartition.LogAppendTime),
					}

					produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, produceResponsePartition)
				}

				response.Responses = append(response.Responses, produceResponse)
			}
		}
	}

	if request.ACKs == 0 {
		return nil, nil
	}

	return response, nil
}
