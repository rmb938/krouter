package v11

import (
	"bytes"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	"github.com/rmb938/krouter/pkg/kafka/records"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
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

		cluster, topic := client.Broker.GetTopic(requestedTopic.Name)

		for _, partition := range requestedTopic.Partitions {
			log = log.WithValues("topic", requestedTopic.Name, "partition", partition.Partition)

			if responseSize >= int64(request.MaxBytes) {
				continue
			}

			partitionResponse := v11.FetchPartitionResponse{
				PartitionIndex:       partition.Partition,
				PreferredReadReplica: 1,
			}

			log = log.WithValues("partition", partition.Partition)

			if topic == nil {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			if partition.Partition >= topic.Partitions {
				partitionResponse.ErrCode = errors.UnknownTopicOrPartition
				topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
				continue
			}

			kafkaFetchRequest := &sarama.FetchRequest{
				MaxWaitTime:  int32(request.MaxWait.Milliseconds()),
				MinBytes:     request.MinBytes,
				MaxBytes:     request.MaxBytes,
				Version:      request.Version(),
				Isolation:    sarama.IsolationLevel(request.IsolationLevel),
				SessionID:    request.SessionID,
				SessionEpoch: request.SessionEpoch,
				RackID:       request.RackID,
			}

			kafkaFetchRequest.AddBlock(topic.Name, partition.Partition, partition.FetchOffset, partition.PartitionMaxBytes)

			kafkaFetchResponse, err := cluster.Fetch(topic, partition.Partition, kafkaFetchRequest)
			if err != nil {
				log.Error(err, "Error fetch to backend cluster")
				if kafkaError, ok := err.(sarama.KError); ok {
					partitionResponse.ErrCode = errors.KafkaError(kafkaError)
				} else {
					return fmt.Errorf("error fetchto backend cluster: %w", err)
				}
			}

			if kafkaFetchResponse != nil {
				if kafkaFetchResponse.ThrottleTime > response.ThrottleDuration {
					response.ThrottleDuration = kafkaFetchResponse.ThrottleTime
				}

				// something returned a high level error so stop everything and just return nothing
				if kafkaFetchResponse.ErrorCode != int16(errors.None) {
					response.ErrCode = errors.KafkaError(kafkaFetchResponse.ErrorCode)
					response.Responses = nil
					break topicLoop
				}

				block := kafkaFetchResponse.Blocks[topicResponse.Topic][partition.Partition]

				if block != nil {
					partitionResponse.ErrCode = errors.KafkaError(block.Err)
					partitionResponse.HighWaterMark = block.HighWaterMarkOffset
					partitionResponse.LastStableOffset = block.LastStableOffset
					partitionResponse.LogStartOffset = block.LogStartOffset

					for _, blockAbortedTransaction := range block.AbortedTransactions {
						partitionResponse.AbortedTransactions = append(partitionResponse.AbortedTransactions, v11.FetchAbortedTransaction{
							ProducerID:  blockAbortedTransaction.ProducerID,
							FirstOffset: blockAbortedTransaction.FirstOffset,
						})
					}

					responseRecordByteBuff := bytes.NewBuffer([]byte{})

					for _, recordBatch := range block.RecordsSet {
						if recordBatch.RecordBatch == nil {
							// record batch is nil, which means we found an old style MsgSet which we can't handle atm
							partitionResponse.ErrCode = errors.InvalidRecord
							continue
						}

						kafkaRb := recordBatch.RecordBatch

						rb := records.RecordBatch{
							BaseOffset:           kafkaRb.FirstOffset,
							PartitionLeaderEpoch: kafkaRb.PartitionLeaderEpoch,
							Attributes:           records.ComputeAttributes(int16(kafkaRb.Codec), kafkaRb.Control, kafkaRb.LogAppendTime, kafkaRb.IsTransactional),
							LastOffsetDelta:      kafkaRb.LastOffsetDelta,
							FirstTimestamp:       kafkaRb.FirstTimestamp,
							MaxTimestamp:         kafkaRb.MaxTimestamp,
							ProducerID:           kafkaRb.ProducerID,
							ProducerEpoch:        kafkaRb.ProducerEpoch,
							BaseSequence:         kafkaRb.FirstSequence,
						}

						for _, kafkaRecord := range kafkaRb.Records {
							record := &records.Record{
								Attributes:     kafkaRecord.Attributes,
								TimeStampDelta: kafkaRecord.TimestampDelta,
								OffsetDelta:    kafkaRecord.OffsetDelta,
								Key:            kafkaRecord.Key,
								Value:          kafkaRecord.Value,
							}

							for _, kafkaHeader := range kafkaRecord.Headers {
								record.Headers = append(record.Headers, records.RecordHeader{
									Key:   kafkaHeader.Key,
									Value: kafkaHeader.Value,
								})
							}

							rb.Records = append(rb.Records, record)
						}

						rbBytes, err := rb.Encode()
						if err != nil {
							return fmt.Errorf("error encoding record batch: %w", err)
						}

						responseRecordByteBuff.Write(rbBytes)
					}

					responseRecordByteBuffLen := responseRecordByteBuff.Len()
					if responseRecordByteBuffLen > 0 {
						responseSize += int64(responseRecordByteBuffLen)
						partitionResponse.Records = responseRecordByteBuff.Bytes()
					}
				}
			}

			topicResponse.Partitions = append(topicResponse.Partitions, partitionResponse)
		}

		response.Responses = append(response.Responses, topicResponse)
	}

	return client.WriteMessage(response, correlationId)
}
