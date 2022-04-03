package logical_broker

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
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
	kafkaClient := c.franzKafkaClient

	type internalTopic struct {
		Name              string
		Partitions        int32
		ReplicationFactor int16
		Config            map[string]*string
	}

	// TODO: configurable replication
	// TODO: do we always need 1 partition or can it be more? We can't change it later though do to out of order messages
	//  so probably safer to keep it at 1
	internalTopics := []internalTopic{
		{
			Name:              InternalTopicTopicPointers,
			Partitions:        1,
			ReplicationFactor: 1,
			Config: map[string]*string{
				// must be compacted topic
				"cleanup.policy": kadm.StringPtr("compact"),
				// compact at least every 24 hours
				"max.compaction.lag.ms": kadm.StringPtr(strconv.Itoa(int((24 * time.Hour).Milliseconds()))),
				// compression because why not
				"compression.type": kadm.StringPtr("gzip"),
			},
		},
		{
			Name:              InternalTopicTopicConfig,
			Partitions:        1,
			ReplicationFactor: 1,
			Config: map[string]*string{
				// must be compacted topic
				"cleanup.policy": kadm.StringPtr("compact"),
				// compact at least every 24 hours
				"max.compaction.lag.ms": kadm.StringPtr(strconv.Itoa(int((24 * time.Hour).Milliseconds()))),
				// compression because why not
				"compression.type": kadm.StringPtr("gzip"),
			},
		},
		{
			Name:              InternalTopicAcls,
			Partitions:        1,
			ReplicationFactor: 1,
			Config: map[string]*string{
				// must be compacted topic
				"cleanup.policy": kadm.StringPtr("compact"),
				// compact at least every 24 hours
				"max.compaction.lag.ms": kadm.StringPtr(strconv.Itoa(int((24 * time.Hour).Milliseconds()))),
				// compression because why not
				"compression.type": kadm.StringPtr("gzip"),
			},
		},
		{
			Name:              InternalTopicBrokers,
			Partitions:        1,
			ReplicationFactor: 1,
			Config: map[string]*string{
				// we only need to keep brokers for 24 hours
				"retention.ms": kadm.StringPtr(strconv.Itoa(int((24 * time.Hour).Milliseconds()))),
				// compression because why not
				"compression.type": kadm.StringPtr("gzip"),
			},
		},
	}

	kafkaAdminClient := kadm.NewClient(kafkaClient)

	for _, topic := range internalTopics {
		// TODO: if topics exist make sure everything matches, otherwise error

		resp, err := kafkaAdminClient.CreateTopics(context.TODO(), topic.Partitions, topic.ReplicationFactor, topic.Config, topic.Name)
		if err != nil {
			return err
		}

		_, err = resp.On(topic.Name, func(response *kadm.CreateTopicResponse) error {
			if response.Err == kerr.TopicAlreadyExists {
				return nil
			}

			return response.Err
		})
		if err != nil {
			return err
		}

		// Produce control record, so we can easily be in sync without using a timer during consuming
		record := kgo.KeySliceRecord([]byte(InternalControlKey), nil)
		record.Topic = topic.Name
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
			c.log.Info("Topic Pointer Kafka Client Closed")
			os.Exit(1)
			return
		}

		fetchErrors := fetches.Errors()
		for _, fetchError := range fetches.Errors() {
			c.log.Error(fetchError.Err, "Error polling fetches for pointer consumer", "partition", fetchError.Partition)
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

			topicName := string(record.Key)
			cluster := record.Value

			if topicName != InternalControlKey {
				if cluster == nil {
					c.topicPointers.Delete(topicName)
				} else {
					c.topicPointers.Store(topicName, string(cluster))
				}
			}

			c.shouldBeSynced(&c.topicPointerSyncOnce, highWaterMarks, lowWaterMarks)
		}
	}
}

func (c *Controller) shouldBeSynced(once *sync.Once, highWaterMarks, lowWaterMarks map[int32]int64) {
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
