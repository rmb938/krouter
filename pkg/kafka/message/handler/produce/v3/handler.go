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

	if request.TransactionalID != nil {
		for _, topicData := range request.TopicData {
			produceResponse := v3.ProduceResponse{
				Name: topicData.Name,
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

	// get cluster from first topic since all topics and partitions in this request should be to the same backend
	//  if our cache has the wrong one some (or all) of our topics will error and the client will refresh metadata
	cluster, _ := broker.GetTopic(request.TopicData[0].Name)

	if cluster == nil {
		for _, topicData := range request.TopicData {
			produceResponse := v3.ProduceResponse{
				Name: topicData.Name,
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

		return response, nil
	}

	var produceTopics []kmsg.ProduceRequestTopic
	for _, topicData := range request.TopicData {
		produceTopic := kmsg.NewProduceRequestTopic()
		produceTopic.Topic = topicData.Name

		for _, partitionData := range topicData.PartitionData {
			produceTopicPartition := kmsg.NewProduceRequestTopicPartition()
			produceTopicPartition.Partition = partitionData.Index
			produceTopicPartition.Records = partitionData.Records

			produceTopic.Partitions = append(produceTopic.Partitions, produceTopicPartition)
		}

		produceTopics = append(produceTopics, produceTopic)
	}

	kafkaResponse, err := cluster.Produce(request.TransactionalID, int32(request.TimeoutDuration.Milliseconds()), produceTopics)
	if err != nil {
		log.Error(err, "Error producing message to backend cluster")
		return nil, fmt.Errorf("error producing to kafka: %w", err)
	}

	if request.ACKs == 0 {
		return nil, nil
	}

	response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond

	for _, kafkaRespTopic := range kafkaResponse.Topics {
		produceResponse := v3.ProduceResponse{
			Name: kafkaRespTopic.Topic,
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

	return response, nil
}
