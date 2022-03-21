package logical_broker

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	TopicConfigRedisKeyFmtPrefix = "{topic-config}-topic-pointer"
	TopicConfigRedisKeyFmt       = TopicConfigRedisKeyFmtPrefix + "-%s"
)

func (c *Controller) APISetTopicPointer(topicName string, cluster *Cluster) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	topicRedisKey := fmt.Sprintf(TopicConfigRedisKeyFmt, topicName)
	err := c.cluster.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {

		exists, err := tx.Exists(ctx, topicRedisKey).Result()
		if err != nil {
			return err
		}
		if exists != 0 {
			return fmt.Errorf("topic %s already exists", topicName)
		}

		_, err = tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			return pipeliner.Set(ctx, topicRedisKey, fmt.Sprintf(TopicConfigClusterRedisKeyFmt, cluster.Name, topicName), 0).Err()
		})

		return err
	}, topicRedisKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) APIDeleteTopicPointer(topicName string) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	topicRedisKey := fmt.Sprintf(TopicConfigRedisKeyFmt, topicName)
	err := c.cluster.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {

		exists, err := tx.Exists(ctx, topicRedisKey).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			return fmt.Errorf("topic %s does not exist", topicName)
		}

		_, err = tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			return pipeliner.Del(ctx, topicRedisKey).Err()
		})

		return err
	}, topicRedisKey)
	if err != nil {
		return err
	}

	return nil
}
