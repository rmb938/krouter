package logical_broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	InternalTopicTopicLeader = "__krouter_topic_leaders"
)

type TopicLeaderMessage struct {
	Name      string `json:"name"`
	Cluster   string `json:"cluster"`
	Partition int32  `json:"partition"`
	Leader    int    `json:"leader"`
}

func (c *Cluster) ConsumeTopicLeaders() {
	kafkaClient, err := c.controller.newFranzKafkaClient(InternalTopicTopicLeader)
	if err != nil {
		c.log.Error(err, "error creating new kafka client for consuming leaders")
		os.Exit(1)
		return
	}
	defer kafkaClient.Close()

	highWaterMarks := make(map[int32]int64)
	lowWaterMarks := make(map[int32]int64)

	for {
		fetches := kafkaClient.PollFetches(context.TODO())
		if fetches.IsClientClosed() {
			c.log.Info("Topic Leader Kafka Client Closed")
			os.Exit(1)
			return
		}

		fetchErrors := fetches.Errors()
		for _, fetchError := range fetches.Errors() {
			c.log.Error(fetchError.Err, "Error polling fetches for leader consumer", "partition", fetchError.Partition)
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
				topicLeaderMessage := &TopicLeaderMessage{}
				err := json.Unmarshal(record.Value, topicLeaderMessage)
				if err != nil {
					// We don't exit and return here because it'll crash all instances
					//  instead we just ignore the message
					c.log.Error(err, "error parsing topic leader", "key", key, "data", string(record.Value))
				}

				if topicLeaderMessage.Cluster == c.Name {
					c.topicLeaderCache.SetWithTTL(fmt.Sprintf(ClusterTopicLeaderKeyFmt, topicLeaderMessage.Name, topicLeaderMessage.Partition), topicLeaderMessage.Leader, 64, 1*time.Hour)
				}
			}

			c.shouldBeSynced(highWaterMarks, lowWaterMarks)
		}
	}
}
