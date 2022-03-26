package logical_broker

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	InternalTopicTopicConfig = "__krouter_topic_configs"
)

type TopicMessageAction string

const (
	TopicMessageActionCreate = "create"
	TopicMessageActionUpdate = "update"
	TopicMessageActionDelete = "delete"
)

type TopicMessage struct {
	Name    string             `json:"name"`
	Action  TopicMessageAction `json:"action"`
	Cluster string             `json:"cluster"`
	Topic   *topics.Topic      `json:"topic"`
}

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
			c.log.Info("Topic Config Kafka Client Closed")
			os.Exit(1)
			return
		}

		fetchErrors := fetches.Errors()
		for _, fetchError := range fetches.Errors() {
			c.log.Error(fetchError.Err, "Error polling fetches for config consumer", "partition", fetchError.Partition)
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

			key := string(record.Key)

			if key != InternalControlKey {
				topicMessage := &TopicMessage{}
				err := json.Unmarshal(record.Value, topicMessage)
				if err != nil {
					// We don't exit and return here because it'll crash all instances
					//  instead we just ignore the message
					c.log.Error(err, "error parsing topic config", "key", key, "data", string(record.Value))
				}

				if topicMessage.Cluster == c.Name {
					if topicMessage.Action == TopicMessageActionDelete {
						c.topics.Delete(topicMessage.Name)
					} else {
						c.topics.Store(topicMessage.Name, topicMessage.Topic)
					}
				}
			}

			c.shouldBeSynced(&c.topicConfigSyncedOnce, highWaterMarks, lowWaterMarks)
		}
	}
}

func (c *Cluster) shouldBeSynced(once *sync.Once, highWaterMarks, lowWaterMarks map[int32]int64) {
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
			c.syncedChan <- struct{}{}
		})
	}
}
