package logical_broker

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/rmb938/krouter/pkg/kafka/logical_broker/internal_topics_pb"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

const (
	InternalTopicBrokers = "__krouter_brokers"
)

func (b *LogicalBroker) ConsumeBrokers() {
	kafkaClient, err := b.controller.newFranzKafkaClient(InternalTopicBrokers)
	if err != nil {
		b.log.Error(err, "error creating new kafka client for consuming brokers")
		os.Exit(1)
		return
	}
	defer kafkaClient.Close()

	highWaterMarks := make(map[int32]int64)
	lowWaterMarks := make(map[int32]int64)

	for {
		fetches := kafkaClient.PollFetches(context.TODO())
		if fetches.IsClientClosed() {
			b.log.Info("Brokers Kafka Client Closed")
			os.Exit(1)
			return
		}

		fetchErrors := fetches.Errors()
		for _, fetchError := range fetches.Errors() {
			b.log.Error(fetchError.Err, "Error polling fetches for brokers consumer", "partition", fetchError.Partition)
		}
		if len(fetchErrors) > 0 {
			os.Exit(1)
		}

		fetches.EachTopic(func(topic kgo.FetchTopic) {
			topic.EachPartition(func(partition kgo.FetchPartition) {
				if _, ok := lowWaterMarks[partition.Partition]; !ok {
					lowWaterMarks[partition.Partition] = 0
				}
				highWaterMarks[partition.Partition] = partition.HighWatermark
			})
		})

		for iter := fetches.RecordIter(); !iter.Done(); {
			record := iter.Next()

			lowWaterMarks[record.Partition] = record.Offset + 1

			key := record.Key
			value := record.Value

			if string(key) != InternalControlKey {
				// only care about messages that are less than 30 seconds old
				if time.Now().Sub(record.Timestamp) < 30*time.Second {
					brokerKey := &internal_topics_pb.BrokerMessageKey{}
					brokerValue := &internal_topics_pb.BrokerMessageValue{}
					err = proto.Unmarshal(key, brokerKey)
					if err != nil {
						// We don't exit and return here because it'll crash all instances
						//  instead we just ignore the message
						b.log.Error(err, "error parsing broker key", "key", key)
						brokerKey = nil
					}

					err = proto.Unmarshal(value, brokerValue)
					if err != nil {
						// We don't exit and return here because it'll crash all instances
						//  instead we just ignore the message
						b.log.Error(err, "error parsing broker value", "value", value)
						brokerValue = nil
					}

					if brokerKey != nil && brokerValue != nil {
						broker := &models.Broker{
							ID: brokerKey.BrokerId,
							Endpoint: models.BrokerEndpoint{
								Host: brokerValue.Endpoint.Host,
								Port: brokerValue.Endpoint.Port,
							},
							Rack:          brokerValue.Rack,
							LastHeartbeat: record.Timestamp,
						}

						b.brokers.Store(strconv.FormatInt(int64(broker.ID), 10), broker)
					}
				}
			}

			b.shouldBeSynced(&b.brokerSyncedOnce, highWaterMarks, lowWaterMarks)
		}
	}
}

func (b *LogicalBroker) shouldBeSynced(once *sync.Once, highWaterMarks, lowWaterMarks map[int32]int64) {
	synced := true
	for partitionIndex, highMark := range highWaterMarks {
		if lowMark, ok := lowWaterMarks[partitionIndex]; ok {
			if highMark != lowMark {
				synced = false
			}
		}
	}

	if synced {
		once.Do(func() {
			b.syncedChan <- struct{}{}
		})
	}
}
