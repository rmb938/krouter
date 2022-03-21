package logical_broker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
)

const (
	TopicConfigClusterRedisKeyFmtPrefix = "{topic-config}-cluster"
	TopicConfigClusterRedisKeyFmtSuffix = "topic-%s"
	TopicConfigClusterRedisKeyFmt       = TopicConfigClusterRedisKeyFmtPrefix + "-%s-" + TopicConfigClusterRedisKeyFmtSuffix
)

func (c *Cluster) APICreateTopic(topicName string, partitions int32, replicationFactor int16, config map[string]*string) (*topics.Topic, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	topicRedisKey := fmt.Sprintf(TopicConfigClusterRedisKeyFmt, c.Name, topicName)
	err := c.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {

		exists, err := tx.Exists(ctx, topicRedisKey).Result()
		if err != nil {
			return err
		}
		if exists != 0 {
			return fmt.Errorf("topic %s already exists", topicName)
		}

		adminClient := kadm.NewClient(c.franzKafkaClient)
		resp, err := adminClient.CreateTopics(ctx, partitions, replicationFactor, config, topicName)
		if err != nil {
			return err
		}

		if err := resp[topicName].Err; err != nil {
			return fmt.Errorf("error from kafka: %w", err)
		}

		hashValues := map[string]interface{}{
			"partitions": partitions,
		}

		for configKey, configValue := range config {
			hashValues[fmt.Sprintf("config.%s", configKey)] = configValue
		}

		_, err = tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			return pipeliner.HSet(ctx, topicRedisKey, hashValues).Err()
		})

		return err
	}, topicRedisKey)
	if err != nil {
		return nil, err
	}

	topic := &topics.Topic{
		Name:       topicName,
		Partitions: partitions,
		Config:     config,
	}

	return topic, nil
}

func (c *Cluster) apiParseTopic(ctx context.Context, tx *redis.Tx, topicName string) (*topics.Topic, error) {
	topicRedisKey := fmt.Sprintf(TopicConfigClusterRedisKeyFmt, c.Name, topicName)
	data, err := tx.HGetAll(ctx, topicRedisKey).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, redis.Nil
	}

	partitions, _ := strconv.Atoi(data["partitions"])
	config := make(map[string]*string)

	for key, value := range data {
		if strings.HasPrefix(key, "config.") {
			configKey := strings.TrimPrefix(key, "config.")
			config[configKey] = &value
		}
	}

	topic := &topics.Topic{
		Name:       topicName,
		Partitions: int32(partitions),
		Config:     config,
	}

	return topic, nil
}

func (c *Cluster) APIGetTopic(topicName string) (*topics.Topic, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var topic *topics.Topic
	topicRedisKey := fmt.Sprintf(TopicConfigClusterRedisKeyFmt, c.Name, topicName)
	err := c.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {
		var err error
		topic, err = c.apiParseTopic(ctx, tx, topicName)
		if err != nil {
			if err == redis.Nil {
				return nil
			}
			return err
		}

		return nil
	}, topicRedisKey)

	return topic, err
}

func (c *Cluster) APIUpdateTopic(topicName string, partitions int32, config map[string]*string) (*topics.Topic, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	topicRedisKey := fmt.Sprintf(TopicConfigRedisKeyFmt, topicName)
	err := c.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {
		topic, err := c.apiParseTopic(ctx, tx, topicName)
		if err != nil {
			if err == redis.Nil {
				return fmt.Errorf("topic %s does not exist", topicName)
			}
			return err
		}

		adminClient := kadm.NewClient(c.franzKafkaClient)
		if topic.Partitions > partitions {
			// update partitions
			resp, err := adminClient.UpdatePartitions(ctx, int(partitions), topicName)
			if err != nil {
				return err
			}
			if err := resp[topicName].Err; err != nil {
				return fmt.Errorf("error from kafka while updating partitions: %w", err)
			}
		} else if topic.Partitions < partitions {
			return fmt.Errorf("cannot decrease topic partitions")
		}

		// calculate config to update
		configToSet := make(map[string]*string)
		configToDelete := make([]string, 0)

		for oldConfigKey, oldConfigValue := range topic.Config {
			if _, ok := config[oldConfigKey]; !ok {
				configToDelete = append(configToDelete, oldConfigKey)
			} else {
				if *oldConfigValue != *config[oldConfigKey] {
					configToSet[oldConfigKey] = oldConfigValue
				}
			}
		}

		if len(configToSet) > 0 || len(configToDelete) > 0 {
			// update config
			var alter []kadm.AlterConfig

			for configKey, configValue := range configToSet {
				alter = append(alter, kadm.AlterConfig{
					Op:    kadm.SetConfig,
					Name:  configKey,
					Value: configValue,
				})
			}

			for _, configKey := range configToDelete {
				alter = append(alter, kadm.AlterConfig{
					Op:    kadm.DeleteConfig,
					Name:  configKey,
					Value: nil,
				})
			}

			resp, err := adminClient.AlterTopicConfigs(ctx, alter, topicName)
			if err != nil {
				return err
			}
			for _, r := range resp {
				if r.Err != nil {
					return fmt.Errorf("error from kafka while updating config: %w", err)
				}
			}
		}

		hashValues := map[string]interface{}{
			"name":       topic.Name,
			"partitions": partitions,
		}

		for configKey, configValue := range config {
			hashValues[fmt.Sprintf("config.%s", configKey)] = configValue
		}

		_, err = tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			return pipeliner.HSet(ctx, topicRedisKey, hashValues).Err()
		})

		return err
	}, topicRedisKey)
	if err != nil {
		return nil, err
	}

	topic := &topics.Topic{
		Name:       topicName,
		Partitions: partitions,
		Config:     config,
	}

	return topic, nil
}

func (c *Cluster) APIDeleteTopicCluster(topicName string) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	topicClusterRedisKey := fmt.Sprintf(TopicConfigClusterRedisKeyFmt, c.Name, topicName)
	err := c.redisClient.Client.Watch(ctx, func(tx *redis.Tx) error {

		exists, err := tx.Exists(ctx, topicClusterRedisKey).Result()
		if err != nil {
			return err
		}
		if exists != 1 {
			return fmt.Errorf("topic %s does not exist", topicName)
		}

		adminClient := kadm.NewClient(c.franzKafkaClient)
		resp, err := adminClient.DeleteTopics(ctx, topicName)
		if err != nil {
			return err
		}
		if err := resp[topicName].Err; err != nil {
			if err != kerr.UnknownTopicOrPartition {
				return fmt.Errorf("error from kafka: %w", err)
			}
		}

		_, err = tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			return pipeliner.Del(ctx, topicClusterRedisKey).Err()
		})

		return err
	}, topicClusterRedisKey)
	if err != nil {
		return err
	}

	return nil
}
