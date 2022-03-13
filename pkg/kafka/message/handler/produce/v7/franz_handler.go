package v7

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/net/message"
)

type FranzHandler struct {
}

func (h *FranzHandler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("franz-produce-v7-handler")
	request := message.(*v7.Request)

	response := &v7.Response{}

	for _, topicData := range request.TopicData {
		log = log.WithValues("topic", topicData.Name)
		produceResponse := v7.ProduceResponse{
			Name: topicData.Name,
		}

		cluster, topic := client.Broker.GetTopic(topicData.Name)

		for _, partitionData := range topicData.PartitionData {
			log = log.WithValues("partition", partitionData.Index)
			partitionResponse := v7.PartitionResponse{
				Index:   partitionData.Index,
				ErrCode: errors.None,
			}

			if topic == nil {
				// Topic is not found so don't do anything else
				log.V(1).Info("Client tried to produce to a topic that doesn't exist")
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
				continue
			}

			if partitionData.Index >= topic.Partitions {
				// Partition is not found so don't do anything else
				log.V(1).Info("Client tried to produce to a topic partition that doesn't exist")
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
				continue
			}

			kafkaResponse, err := cluster.FranzProduce(topic, partitionData.Index, request.TransactionalID, int32(request.TimeoutDuration.Milliseconds()), partitionData.Records)
			if err != nil {
				log.Error(err, "Error producing message to backend cluster")
				return fmt.Errorf("error producing to kafka: %w", err)
			}

			if int64(kafkaResponse.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
				response.ThrottleDuration = time.Duration(kafkaResponse.ThrottleMillis) * time.Millisecond
			}

			partitionResp := kafkaResponse.Topics[0].Partitions[0]

			partitionResponse.ErrCode = errors.KafkaError(partitionResp.ErrorCode)
			partitionResponse.BaseOffset = partitionResp.BaseOffset
			partitionResponse.LogStartOffset = partitionResp.LogStartOffset
			partitionResponse.LogAppendTime = time.UnixMilli(partitionResp.LogAppendTime)

			produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
		}
		response.Responses = append(response.Responses, produceResponse)
	}

	return client.WriteMessage(response, correlationId)
}
