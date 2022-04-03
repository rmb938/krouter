package v1

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v1"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("produce-v1-handler")
	request := message.(*v1.Request)

	response := &v1.Response{}

	if len(request.TopicData) == 0 {
		return response, nil
	}

	response.Responses = make([]v1.ProduceResponse, 0, len(request.TopicData))

	if request.ACKs < -1 || request.ACKs > 1 {
		for _, topicData := range request.TopicData {
			produceResponse := v1.ProduceResponse{
				Name:               topicData.Name,
				PartitionResponses: make([]v1.PartitionResponse, 0, len(topicData.PartitionData)),
			}
			for _, partitionData := range topicData.PartitionData {
				partitionResponse := v1.PartitionResponse{
					Index:   partitionData.Index,
					ErrCode: errors.InvalidRequiredAcks,
				}
				produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
			}
			response.Responses = append(response.Responses, produceResponse)
		}

		return response, nil
	}

	clusterRequests := make(map[*logical_broker.Cluster][]kmsg.ProduceRequestTopic)
	for _, topicData := range request.TopicData {
		cluster := broker.GetClusterByTopic(topicData.Name)
		if cluster == nil {
			produceResponse := v1.ProduceResponse{
				Name:               topicData.Name,
				PartitionResponses: make([]v1.PartitionResponse, 0, len(topicData.PartitionData)),
			}
			for _, partitionData := range topicData.PartitionData {
				partitionResponse := v1.PartitionResponse{
					Index:   partitionData.Index,
					ErrCode: errors.UnknownTopicOrPartition,
				}
				produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
			}
			response.Responses = append(response.Responses, produceResponse)
			continue
		}

		clusterRequest, ok := clusterRequests[cluster]
		if !ok {
			clusterRequest = make([]kmsg.ProduceRequestTopic, 0)
		}

		produceRequestTopic := kmsg.NewProduceRequestTopic()
		produceRequestTopic.Topic = topicData.Name

		for _, partitionData := range topicData.PartitionData {
			produceRequestPartition := kmsg.NewProduceRequestTopicPartition()
			produceRequestPartition.Partition = partitionData.Index
			produceRequestPartition.Records = partitionData.Records

			produceRequestTopic.Partitions = append(produceRequestTopic.Partitions, produceRequestPartition)
		}

		clusterRequest = append(clusterRequest, produceRequestTopic)
		clusterRequests[cluster] = clusterRequest
	}

	for cluster, produceRequests := range clusterRequests {
		kafkaResponse, err := cluster.Produce(broker.BrokerID, nil, int32(request.TimeoutDuration.Milliseconds()), produceRequests)
		if err != nil {
			log.Error(err, "Error producing message to backend cluster")
			return nil, fmt.Errorf("error producing to kafka: %w", err)
		}

		if request.ACKs != 0 {
			if kafkaResponse.ThrottleMillis > int32(response.ThrottleDuration) {
				response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond
			}

			for _, kafkaRespTopic := range kafkaResponse.Topics {
				produceResponse := v1.ProduceResponse{}
				produceResponse.Name = kafkaRespTopic.Topic

				for _, partitionResponse := range kafkaRespTopic.Partitions {
					produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, v1.PartitionResponse{
						Index:      partitionResponse.Partition,
						ErrCode:    errors.KafkaError(partitionResponse.ErrorCode),
						BaseOffset: partitionResponse.BaseOffset,
					})
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
