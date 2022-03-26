package logical_broker

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	InternalControlKey         = "__krouter_control"
	InternalTopicTopicPointers = "__krouter_topic_pointers"
)

func (c *Controller) CreateInternalTopics() error {
	kafkaClient, err := c.newFranzKafkaClient()
	if err != nil {
		return err
	}
	defer kafkaClient.Close()

	kafkaAdminClient := kadm.NewClient(kafkaClient)

	for _, topicName := range []string{InternalTopicTopicPointers, InternalTopicTopicConfig} {
		// TODO: configurable partitions and replication
		resp, err := kafkaAdminClient.CreateTopics(context.TODO(), 5, 1, map[string]*string{
			// must be compacted topic
			"cleanup.policy": kadm.StringPtr("compact"),
			// compact at least every 24 hours
			"max.compaction.lag.ms": kadm.StringPtr(strconv.Itoa(int((24 * time.Hour).Milliseconds()))),
			// compression because why not
			"compression.type": kadm.StringPtr("gzip"),
		}, topicName)
		if err != nil {
			return err
		}

		_, err = resp.On(topicName, func(response *kadm.CreateTopicResponse) error {
			if response.Err == kerr.TopicAlreadyExists {
				return nil
			}

			return response.Err
		})
		if err != nil {
			return err
		}

		record := kgo.KeyStringRecord(InternalControlKey, "")
		record.Topic = topicName
		produceResp := kafkaClient.ProduceSync(context.TODO(), record)
		if produceResp.FirstErr() != nil {
			return fmt.Errorf("error control message to internal topics: %w", produceResp.FirstErr())
		}
	}

	return nil
}

func (c *Controller) ConsumeTopicPointers() {
	kafkaClient, err := c.newFranzKafkaClient(InternalTopicTopicPointers)
	if err != nil {
		c.log.Error(err, "error creating new kafka client for consuming pointers")
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
			cluster := string(record.Value)

			if topicName != InternalControlKey {
				if len(cluster) == 0 {
					c.topicPointers.Delete(topicName)
				} else {
					c.topicPointers.Store(topicName, cluster)
				}
			}

			c.shouldBeSynced(highWaterMarks, lowWaterMarks)
		}
	}
}

func (c *Controller) shouldBeSynced(highWaterMarks, lowWaterMarks map[int32]int64) {
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
			c.log.Info("Controller Synced Topic Pointers")
			c.syncedChan <- struct{}{}
		})
	}
}
