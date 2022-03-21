package v7

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/net/message"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("produce-v7-handler")
	request := message.(*v7.Request)

	response := &v7.Response{}

	if len(request.TopicData) == 0 {
		return response, nil
	}

	// get cluster from first topic since all topics and partitions in this request should be to the same backend
	//  if our cache has the wrong one some (or all) of our topics will error and the client will refresh metadata
	cluster, _, err := broker.GetTopic(request.TopicData[0].Name)
	if err != nil {
		log.Error(err, "error getting topic from logical broker")
		return nil, fmt.Errorf("error getting topic from logical broker: %w", err)
	}

	if cluster == nil {
		for _, topicData := range request.TopicData {
			produceResponse := v7.ProduceResponse{
				Name: topicData.Name,
			}
			for _, partitionData := range topicData.PartitionData {
				partitionResponse := v7.PartitionResponse{
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

	response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond

	for _, kafkaRespTopic := range kafkaResponse.Topics {
		produceResponse := v7.ProduceResponse{
			Name: kafkaRespTopic.Topic,
		}

		for _, kafkaRespTopicPartition := range kafkaRespTopic.Partitions {
			produceResponsePartition := v7.PartitionResponse{
				Index:          kafkaRespTopicPartition.Partition,
				ErrCode:        errors.KafkaError(kafkaRespTopicPartition.ErrorCode),
				BaseOffset:     kafkaRespTopicPartition.BaseOffset,
				LogAppendTime:  time.UnixMilli(kafkaRespTopicPartition.LogAppendTime),
				LogStartOffset: kafkaRespTopicPartition.LogStartOffset,
			}

			produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, produceResponsePartition)
		}

		response.Responses = append(response.Responses, produceResponse)
	}

	// for _, topicData := range request.TopicData {
	// 	log = log.WithValues("topic", topicData.Name)
	// 	produceResponse := v7.ProduceResponse{
	// 		Name: topicData.Name,
	// 	}
	//
	// 	cluster, topic := broker.GetTopic(topicData.Name)
	//
	// 	for _, partitionData := range topicData.PartitionData {
	// 		log = log.WithValues("partition", partitionData.Index)
	// 		partitionResponse := v7.PartitionResponse{
	// 			Index:   partitionData.Index,
	// 			ErrCode: errors.None,
	// 		}
	//
	// 		if topic == nil {
	// 			// Topic is not found so don't do anything else
	// 			log.V(1).Info("Client tried to produce to a topic that doesn't exist")
	// 			partitionResponse.ErrCode = errors.UnknownTopicOrPartition
	// 			produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
	// 			continue
	// 		}
	//
	// 		if partitionData.Index >= topic.Partitions {
	// 			// Partition is not found so don't do anything else
	// 			log.V(1).Info("Client tried to produce to a topic partition that doesn't exist")
	// 			partitionResponse.ErrCode = errors.UnknownTopicOrPartition
	// 			produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
	// 			continue
	// 		}
	//
	// 		kafkaResponse, err := cluster.Produce(topic, partitionData.Index, request.TransactionalID, int32(request.TimeoutDuration.Milliseconds()), partitionData.Records)
	// 		if err != nil {
	// 			log.Error(err, "Error producing message to backend cluster")
	// 			return nil, fmt.Errorf("error producing to kafka: %w", err)
	// 		}
	//
	// 		if int64(kafkaResponse.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
	// 			response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond
	// 		}
	//
	// 		partitionResp := kafkaResponse.Topics[0].Partitions[0]
	//
	// 		partitionResponse.ErrCode = errors.KafkaError(partitionResp.ErrorCode)
	// 		partitionResponse.BaseOffset = partitionResp.BaseOffset
	// 		partitionResponse.LogStartOffset = partitionResp.LogStartOffset
	// 		partitionResponse.LogAppendTime = time.UnixMilli(partitionResp.LogAppendTime)
	//
	// 		produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
	// 	}
	// 	response.Responses = append(response.Responses, produceResponse)
	// }

	return response, nil
}
