package v5

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v5 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v5"
	"github.com/rmb938/krouter/pkg/kafka/records"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, message message.Message, correlationId int32) error {
	request := message.(*v5.Request)

	response := &v5.Response{}

	kafkaClient, err := sarama.NewClient([]string{"localhost:9094"}, sarama.NewConfig())
	if err != nil {
		return fmt.Errorf("error creating sarama client: %w", err)
	}

	defer kafkaClient.Close()

	for _, topicData := range request.TopicData {
		produceResponse := v5.ProduceResponse{
			Name: topicData.Name,
		}

		for _, partitionData := range topicData.PartitionData {
			partitionResponse := v5.PartitionResponse{
				Index:   partitionData.Index,
				ErrCode: errors.None,
			}

			rb, err := records.ParseRecordBatch(partitionData.Records)
			if err != nil {
				return fmt.Errorf("error parsing record patch: %w", err)
			}

			broker, err := kafkaClient.Leader(topicData.Name, partitionData.Index)
			if err != nil {
				if kafkaError, ok := err.(sarama.KError); ok {
					partitionResponse.ErrCode = errors.KafkaError(kafkaError)
					produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
					continue
				}

				return fmt.Errorf("error finding kafka topic partition leader: %w", err)
			}
			kafkaProduceRequest := &sarama.ProduceRequest{
				TransactionalID: request.TransactionalID,
				RequiredAcks:    sarama.RequiredAcks(request.ACKs),
				Timeout:         int32(request.TimeoutDuration.Milliseconds()),
				Version:         request.Version(),
			}

			kafkaRb := &sarama.RecordBatch{
				FirstOffset:           rb.BaseOffset,
				PartitionLeaderEpoch:  rb.PartitionLeaderEpoch,
				Version:               rb.Magic,
				Codec:                 0,     // TODO: pull from rb attributes
				CompressionLevel:      0,     // TODO: pull from rb attributes
				Control:               false, // TODO: pull from rb attributes
				LogAppendTime:         false, // TODO: pull from rb attributes
				LastOffsetDelta:       rb.LastOffsetDelta,
				FirstTimestamp:        rb.FirstTimestamp,
				MaxTimestamp:          rb.MaxTimestamp,
				ProducerID:            rb.ProducerID,
				ProducerEpoch:         rb.ProducerEpoch,
				FirstSequence:         rb.BaseSequence,
				PartialTrailingRecord: false, // TODO: pull from rb attributes
				IsTransactional:       false, // TODO: pull from rb attributes
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

			kafkaProduceResponse, err := broker.Produce(kafkaProduceRequest)
			if err != nil {
				if kafkaError, ok := err.(sarama.KError); ok {
					partitionResponse.ErrCode = errors.KafkaError(kafkaError)
					produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
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
