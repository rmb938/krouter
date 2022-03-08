package logical_broker

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type Controller struct {
	redisClient redis.UniversalClient
	cluster     *Cluster
}

func NewController(log logr.Logger, cluster *Cluster, redisAddresses []string) (*Controller, error) {
	log = log.WithName("controller")

	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: redisAddresses,
	})

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer redisContextCancel()

	_, err := redisClient.Ping(redisContext).Result()
	if err != nil {
		return nil, fmt.Errorf("redis: error pining redis in controller: %w", err)
	}

	appendonly, err := redisClient.ConfigGet(redisContext, "appendonly").Result()
	if err != nil {
		return nil, fmt.Errorf("redis: error checking appendonly config: %w", err)
	}

	if appendonly[1] != "yes" {
		return nil, fmt.Errorf("redis: appendonly config must be set to 'yes': Value: '%v'", appendonly[1])
	}

	appendfsync, err := redisClient.ConfigGet(redisContext, "appendfsync").Result()
	if err != nil {
		return nil, fmt.Errorf("redis: error checking appendfsync config: %w", err)
	}
	if appendfsync[1] != "always" && appendfsync[1] != "everysec" {
		return nil, fmt.Errorf("redis: appendfsync config must be set to 'always' or 'everysec': Value: '%v'", appendonly[1])
	}

	return &Controller{
			redisClient: redisClient,
			cluster:     cluster},
		nil
}

func (c *Controller) findCoordinator(consumerGroup string) (*sarama.Broker, error) {
	return c.cluster.kafkaClient.Coordinator(consumerGroup)
}

func (c *Controller) JoinGroup(request *sarama.JoinGroupRequest) (*sarama.JoinGroupResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.GroupId)
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
	err = c.redisClient.Watch(redisContext, func(tx *redis.Tx) error {
		// TODO: make exp configurable
		return tx.Set(redisContext, redisGroupGenerationKey, response.GenerationId, 7*24*time.Hour).Err()
	}, redisGroupGenerationKey)

	return response, err
}

func (c *Controller) SyncGroup(request *sarama.SyncGroupRequest) (*sarama.SyncGroupResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.GroupId)
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.SyncGroup(request)
}

func (c *Controller) LeaveGroup(request *sarama.LeaveGroupRequest) (*sarama.LeaveGroupResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.GroupId)
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.LeaveGroup(request)
}

func (c *Controller) HeartBeat(request *sarama.HeartbeatRequest) (*sarama.HeartbeatResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.GroupId)
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
	err = c.redisClient.Watch(redisContext, func(tx *redis.Tx) error {
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

func (c *Controller) OffsetFetch(group, topic string, partition int32) (int64, error) {
	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	redisGroupOffsetKey := fmt.Sprintf("{group-%s}-offset-topic-%s-partition-%d", group, topic, partition)

	var offset int64
	err := c.redisClient.Watch(redisContext, func(tx *redis.Tx) error {
		var err error
		offset, err = tx.Get(redisContext, redisGroupOffsetKey).Int64()

		if err == redis.Nil {
			return nil
		}

		return err
	}, redisGroupOffsetKey)

	return offset, err
}

func (c *Controller) OffsetCommit(group, topic string, groupGenerationId, partition int32, offset int64) error {
	// TODO: to expire these, every 5 minutes do a `SCAN MATCH group-offset-*` and see if a generation exists, if it doesn't delete it

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer redisContextCancel()

	redisGroupGenerationKey := fmt.Sprintf("{group-%s}-generation", group)
	redisGroupOffsetKey := fmt.Sprintf("{group-%s}-offset-topic-%s-partition-%d", group, topic, partition)

	err := c.redisClient.Watch(redisContext, func(tx *redis.Tx) error {
		generationId, err := tx.Get(redisContext, redisGroupGenerationKey).Int64()
		if err != nil && err != redis.Nil {
			return err
		}

		if err == redis.Nil || generationId != int64(groupGenerationId) {
			return sarama.KError(errors.IllegalGeneration)
		}

		return tx.Set(redisContext, redisGroupOffsetKey, offset, 0).Err()
	}, redisGroupGenerationKey)

	return err
}
