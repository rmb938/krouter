package v7

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/records"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("produce-v7-handler")
	request := message.(*v7.Request)

	response := &v7.Response{}

	// TODO: move this client to a central location
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_6_0_0
	kafkaClient, err := sarama.NewClient([]string{"localhost:9094"}, saramaConfig)
	if err != nil {
		return fmt.Errorf("error creating sarama client: %w", err)
	}

	defer kafkaClient.Close()

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

			rb, err := records.ParseRecordBatch(partitionData.Records)
			if err != nil {
				return fmt.Errorf("error parsing record patch: %w", err)
			}

			kafkaProduceRequest := &sarama.ProduceRequest{
				TransactionalID: request.TransactionalID,
				RequiredAcks:    sarama.RequiredAcks(request.ACKs),
				Timeout:         int32(request.TimeoutDuration.Milliseconds()),
				Version:         request.Version(),
			}

			kafkaRb := &sarama.RecordBatch{
				FirstOffset:          rb.BaseOffset,
				PartitionLeaderEpoch: rb.PartitionLeaderEpoch,
				Version:              rb.Magic,
				Codec:                sarama.CompressionCodec(rb.GetCodec()),
				CompressionLevel:     sarama.CompressionLevelDefault, // need to set this so records are re-compressed correctly
				Control:              rb.IsControl(),
				LogAppendTime:        rb.IsLogAppendTime(),
				LastOffsetDelta:      rb.LastOffsetDelta,
				FirstTimestamp:       rb.FirstTimestamp,
				MaxTimestamp:         rb.MaxTimestamp,
				ProducerID:           rb.ProducerID,
				ProducerEpoch:        rb.ProducerEpoch,
				FirstSequence:        rb.BaseSequence,
				IsTransactional:      rb.IsTransactional(),
			}

			for _, r := range rb.Records {
				kafkaRecord := &sarama.Record{
					Attributes:     r.Attributes,
					TimestampDelta: r.TimeStampDelta,
					OffsetDelta:    r.OffsetDelta,
					Key:            r.Key,
					Value:          r.Value,
				}

				for _, rHeader := range r.Headers {
					kafkaRecord.Headers = append(kafkaRecord.Headers, &sarama.RecordHeader{
						Key:   rHeader.Key,
						Value: rHeader.Value,
					})
				}

				kafkaRb.Records = append(kafkaRb.Records, kafkaRecord)
			}

			kafkaProduceRequest.AddBatch(topicData.Name, partitionData.Index, kafkaRb)

			kafkaProduceResponse, err := cluster.Produce(topic, partitionData.Index, kafkaProduceRequest)
			if err != nil {
				log.Error(err, "Error producing message to backend cluster")
				if kafkaError, ok := err.(sarama.KError); ok {
					partitionResponse.ErrCode = errors.KafkaError(kafkaError)
					produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)

					if kafkaProduceResponse.ThrottleTime > response.ThrottleDuration {
						response.ThrottleDuration = kafkaProduceResponse.ThrottleTime
					}
					continue
				}
				return fmt.Errorf("error producing to kafka: %w", err)
			}

			if kafkaProduceResponse.ThrottleTime > response.ThrottleDuration {
				response.ThrottleDuration = kafkaProduceResponse.ThrottleTime
			}

			responseBlock := kafkaProduceResponse.GetBlock(topicData.Name, partitionData.Index)

			partitionResponse.ErrCode = errors.KafkaError(responseBlock.Err)
			partitionResponse.BaseOffset = responseBlock.Offset
			partitionResponse.LogStartOffset = responseBlock.StartOffset
			partitionResponse.LogAppendTime = responseBlock.Timestamp

			produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
		}
		response.Responses = append(response.Responses, produceResponse)
	}

	return client.WriteMessage(response, correlationId)
}
