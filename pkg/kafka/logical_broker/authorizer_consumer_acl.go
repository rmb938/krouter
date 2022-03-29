package logical_broker

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/internal_topics_pb"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

const (
	InternalTopicAcls = "__krouter_acls"
)

func (c *Authorizer) ConsumeAcls() {
	kafkaClient := c.kafkaClient

	highWaterMarks := make(map[int32]int64)
	lowWaterMarks := make(map[int32]int64)

	for {
		fetches := kafkaClient.PollFetches(context.TODO())
		if fetches.IsClientClosed() {
			c.log.Info("Topic ACL Kafka Client Closed")
			os.Exit(1)
			return
		}

		fetchErrors := fetches.Errors()
		for _, fetchError := range fetches.Errors() {
			c.log.Error(fetchError.Err, "Error polling fetches for acl consumer", "partition", fetchError.Partition)
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
				aclKey := &internal_topics_pb.ACLMessageKey{}
				aclValue := &internal_topics_pb.ACLMessageValue{}
				err := proto.Unmarshal(key, aclKey)
				if err != nil {
					// We don't exit and return here because it'll crash all instances
					//  instead we just ignore the message
					c.log.Error(err, "error parsing acl key", "key", key)
					aclKey = nil
				}

				// value may be nil so don't parse if it is
				if value != nil {
					err = proto.Unmarshal(value, aclValue)
					if err != nil {
						// We don't exit and return here because it'll crash all instances
						//  instead we just ignore the message
						c.log.Error(err, "error parsing acl value", "value", value)
						aclValue = nil
					}
				}

				// if the key and value was parsed correctly do stuff
				if aclKey != nil && aclValue != nil {
					aclModel := &models.ACL{
						Operation:    models.ACLOperation(aclKey.Operation),
						ResourceType: models.ACLResourceType(aclKey.ResourceType),
						PatternType:  models.ACLPatternType(aclKey.PatternType),
						ResourceName: aclKey.ResourceName,
						Principal:    aclKey.Principal,
					}

					hashInt, err := hashstructure.Hash(aclModel, hashstructure.FormatV2, nil)
					if err != nil {
						c.log.Error(err, "hashing acl for map key")
					}
					hashKey := strconv.FormatUint(hashInt, 10)

					aclModel.Permission = models.ACLPermission(aclValue.Permission)

					if value == nil {
						// if no value delete it
						c.acls.Delete(hashKey)
					} else {
						c.acls.Store(hashKey, aclModel)
					}
				}
			}

			c.shouldBeSynced(&c.aclsSyncOnce, highWaterMarks, lowWaterMarks)
		}
	}
}

func (c *Authorizer) shouldBeSynced(once *sync.Once, highWaterMarks, lowWaterMarks map[int32]int64) {
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
