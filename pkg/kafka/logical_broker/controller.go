package logical_broker

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/twmb/franz-go/pkg/kmsg"
)

const (
	GroupCoordinatorRedisKeyFmt = "{group-%s}-coordinator"
)

type Controller struct {
	cluster *Cluster
}

func NewController(log logr.Logger, cluster *Cluster) (*Controller, error) {
	log = log.WithName("controller")

	return &Controller{cluster: cluster}, nil
}

func (c *Controller) ClusterMetadata(ctx context.Context) (*kmsg.MetadataResponse, error) {
	return c.cluster.TopicMetadata(ctx, []string{})
}

func (c *Controller) FindCoordinator(consumerGroup string) (*kmsg.FindCoordinatorResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	request := kmsg.NewPtrFindCoordinatorRequest()
	request.CoordinatorType = 0
	request.CoordinatorKey = consumerGroup

	response, err := c.cluster.franzKafkaClient.Request(ctx, request)
	if err != nil {
		return nil, err
	}

	coordinatorResponse := response.(*kmsg.FindCoordinatorResponse)

	if coordinatorResponse.ErrorCode == int16(errors.None) {
		// Consumers don't actively refresh this
		// so if this expires we should return errors.NotCoordinator and clients will try and FindCoordinator again
		err := c.cluster.redisClient.Client.Set(ctx, fmt.Sprintf(GroupCoordinatorRedisKeyFmt, consumerGroup), coordinatorResponse.NodeID, 1*time.Hour).Err()
		if err != nil {
			return nil, err
		}
	}

	return coordinatorResponse, nil
}

func (c *Controller) groupCoordinator(ctx context.Context, consumerGroup string) (int, error) {
	brokerID, err := c.cluster.redisClient.Client.Get(ctx, fmt.Sprintf(GroupCoordinatorRedisKeyFmt, consumerGroup)).Int()
	if err != nil {
		if err == redis.Nil {
			return -1, nil
		}
		return -1, err
	}

	return brokerID, nil
}

func (c *Controller) JoinGroup(request *sarama.JoinGroupRequest) (*sarama.JoinGroupResponse, error) {
	coordinatorId, err := c.groupCoordinator(context.TODO(), request.GroupId)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		return &sarama.JoinGroupResponse{
			Err: sarama.KError(errors.NotCoordinator),
		}, nil
	}

	coordinatorBroker, err := c.cluster.saramaKafkaClient.Broker(int32(coordinatorId))
	if err != nil {
		return nil, err
	}

	response, err := coordinatorBroker.JoinGroup(request)
	if err != nil {
		return response, err
	}

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	redisGroupGenerationKey := fmt.Sprintf("{group-%s}-generation", request.GroupId)
	err = c.cluster.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		// TODO: make exp configurable
		return tx.Set(redisContext, redisGroupGenerationKey, response.GenerationId, 7*24*time.Hour).Err()
	}, redisGroupGenerationKey)

	return response, err
}

func (c *Controller) SyncGroup(request *sarama.SyncGroupRequest) (*sarama.SyncGroupResponse, error) {
	coordinatorId, err := c.groupCoordinator(context.TODO(), request.GroupId)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		return &sarama.SyncGroupResponse{
			Err: sarama.KError(errors.NotCoordinator),
		}, nil
	}

	coordinatorBroker, err := c.cluster.saramaKafkaClient.Broker(int32(coordinatorId))
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.SyncGroup(request)
}

func (c *Controller) LeaveGroup(request *sarama.LeaveGroupRequest) (*sarama.LeaveGroupResponse, error) {
	coordinatorId, err := c.groupCoordinator(context.TODO(), request.GroupId)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		return &sarama.LeaveGroupResponse{
			Err: sarama.KError(errors.NotCoordinator),
		}, nil
	}

	coordinatorBroker, err := c.cluster.saramaKafkaClient.Broker(int32(coordinatorId))
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.LeaveGroup(request)
}

func (c *Controller) HeartBeat(request *sarama.HeartbeatRequest) (*sarama.HeartbeatResponse, error) {
	coordinatorId, err := c.groupCoordinator(context.TODO(), request.GroupId)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		return &sarama.HeartbeatResponse{
			Err: sarama.KError(errors.NotCoordinator),
		}, nil
	}

	coordinatorBroker, err := c.cluster.saramaKafkaClient.Broker(int32(coordinatorId))
	if err != nil {
		return nil, err
	}

	response, err := coordinatorBroker.Heartbeat(request)
	if err != nil {
		return response, err
	}

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	redisGroupGenerationKey := fmt.Sprintf("{group-%s}-generation", request.GroupId)
	err = c.cluster.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		generationId, err := tx.Get(redisContext, redisGroupGenerationKey).Int64()
		if err != nil && err != redis.Nil {
			return err
		}

		if err == redis.Nil || generationId != int64(request.GenerationId) {
			return sarama.KError(errors.IllegalGeneration)
		}

		// TODO: make exp configurable
		return tx.Set(redisContext, redisGroupGenerationKey, request.GenerationId, 7*24*time.Hour).Err()
	}, redisGroupGenerationKey)

	return response, err
}

func (c *Controller) base64Topic(topic string) string {
	return base64.StdEncoding.EncodeToString([]byte(topic))
}

func (c *Controller) OffsetFetch(group, topic string, partition int32) (int64, error) {
	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	redisGroupOffsetKey := fmt.Sprintf("{group-%s}-offset-topic-%s-partition-%d", group, c.base64Topic(topic), partition)

	var offset int64
	err := c.cluster.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		var err error
		offset, err = tx.Get(redisContext, redisGroupOffsetKey).Int64()

		if err == redis.Nil {
			return nil
		}

		return err
	}, redisGroupOffsetKey)

	return offset, err
}

func (c *Controller) OffsetFetchAllTopics(group string) (map[string]map[int32]int64, error) {
	offsets := make(map[string]map[int32]int64)

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer redisContextCancel()

	redisGroupOffsetKeyPrefix := fmt.Sprintf("{group-%s}-offset-topic-", group)

	err := c.cluster.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		scan := tx.Scan(redisContext, 0, fmt.Sprintf("%s*", redisGroupOffsetKeyPrefix), 0).Iterator()

		for scan.Next(redisContext) {
			key := scan.Val()
			offset, err := tx.Get(redisContext, key).Int64()

			if err == redis.Nil {
				continue
			}

			if err != nil {
				return err
			}

			topicParts := strings.Split(strings.TrimPrefix(key, redisGroupOffsetKeyPrefix), "-partition-")

			topicByte, err := base64.StdEncoding.DecodeString(topicParts[0])
			if err != nil {
				return fmt.Errorf("err decoding base64 topic name: %s: %w", topicParts[0], err)
			}
			topic := string(topicByte)

			partition, err := strconv.Atoi(topicParts[1])
			if err != nil {
				return err
			}

			if offsets[topic] == nil {
				offsets[topic] = make(map[int32]int64)
			}
			offsets[topic][int32(partition)] = offset
		}

		return nil
	})

	return offsets, err
}

func (c *Controller) OffsetCommit(group, topic string, groupGenerationId, partition int32, offset int64) error {
	// TODO: to expire these, every 5 minutes do a `SCAN MATCH group-offset-*` and see if a generation exists, if it doesn't delete it

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	redisGroupGenerationKey := fmt.Sprintf("{group-%s}-generation", group)
	redisGroupOffsetKey := fmt.Sprintf("{group-%s}-offset-topic-%s-partition-%d", group, c.base64Topic(topic), partition)

	err := c.cluster.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		// generation will be -1 when we are resetting offsets
		if groupGenerationId != -1 {
			generationId, err := tx.Get(redisContext, redisGroupGenerationKey).Int64()
			if err != nil && err != redis.Nil {
				return err
			}

			if err == redis.Nil || generationId != int64(groupGenerationId) {
				return sarama.KError(errors.IllegalGeneration)
			}
		}

		return tx.Set(redisContext, redisGroupOffsetKey, offset, 0).Err()
	}, redisGroupGenerationKey)

	return err
}

func (c *Controller) DescribeGroup(group string) (*sarama.DescribeGroupsResponse, error) {
	coordinatorId, err := c.groupCoordinator(context.TODO(), group)
	if err != nil {
		return nil, err
	}

	if coordinatorId == -1 {
		return &sarama.DescribeGroupsResponse{
			Groups: []*sarama.GroupDescription{
				{
					Err:     sarama.KError(errors.NotCoordinator),
					GroupId: group,
				},
			},
		}, nil
	}

	coordinatorBroker, err := c.cluster.saramaKafkaClient.Broker(int32(coordinatorId))
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.DescribeGroups(&sarama.DescribeGroupsRequest{Groups: []string{group}})
}
