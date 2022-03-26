package logical_broker

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	InternalTopicTopicConfig = "__krouter_topic_config"
)

func (c *Cluster) ConsumeTopicConfigs() {
	kafkaClient, err := c.controller.newFranzKafkaClient(InternalTopicTopicConfig)
	if err != nil {
		c.log.Error(err, "error creating new kafka client for consuming configs")
		os.Exit(1)
		return
	}
	defer kafkaClient.Close()

	highWaterMarks := make(map[int32]int64)
	lowWaterMarks := make(map[int32]int64)

	for {
		fetches := kafkaClient.PollFetches(context.TODO())
		if fetches.IsClientClosed() {
			c.log.Info("Controller Kafka Client Closed")
			return
		}

		fetches.EachError(func(s string, i int32, err error) {
			c.log.Error(err, "Error polling fetches", "partition", i)
			return
		})

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

			topicName := string(record.Key)

			if topicName != InternalControlKey {
				if record.Value == nil {
					c.topics.Delete(topicName)
				} else {
					topic := &topics.Topic{}
					err := json.Unmarshal(record.Value, topic)
					if err != nil {
						c.log.Error(err, "error parsing topic config", "topic", topicName, "config", string(record.Value))
						os.Exit(1)
						return
					}

					c.topics.Store(topicName, topic)
				}
			}

			c.shouldBeSynced(highWaterMarks, lowWaterMarks)
		}
	}
}

func (c *Cluster) shouldBeSynced(highWaterMarks, lowWaterMarks map[int32]int64) {
	synced := true
	for partitionIndex, highMark := range highWaterMarks {
		if lowMark, ok := lowWaterMarks[partitionIndex]; ok {
			if highMark != lowMark {
				synced = false
			}
		}
	}

	if synced {
		c.syncedOnce.Do(func() {
			c.log.Info("Cluster Synced Topic Configs")
			c.syncedChan <- struct{}{}
		})
	}
}
