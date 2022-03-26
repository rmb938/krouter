package logical_broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (c *Cluster) APICreateTopic(topicName string, partitions int32, replicationFactor int16, config map[string]*string) (*topics.Topic, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	adminClient := kadm.NewClient(c.franzKafkaClient)
	resp, err := adminClient.CreateTopics(ctx, partitions, replicationFactor, config, topicName)
	if err != nil {
		return nil, err
	}

	if err := resp[topicName].Err; err != nil {
		return nil, fmt.Errorf("error from kafka: %w", err)
	}

	topic := &topics.Topic{
		Name:       topicName,
		Partitions: partitions,
		Config:     config,
	}

	topicMessage := &TopicMessage{
		Name:    topicName,
		Action:  TopicMessageActionCreate,
		Cluster: c.Name,
		Topic:   topic,
	}

	kafkaClient := c.controller.franzKafkaClient
	topicMessageBytes, _ := json.Marshal(topicMessage)
	record := kgo.KeySliceRecord([]byte(fmt.Sprintf("cluster-%s-topic-%s", c.Name, topicName)), topicMessageBytes)
	record.Topic = InternalTopicTopicConfig
	produceResp := kafkaClient.ProduceSync(context.TODO(), record)
	if produceResp.FirstErr() != nil {
		return nil, produceResp.FirstErr()
	}

	return topic, nil
}

func (c *Cluster) APIGetTopic(topicName string) (*topics.Topic, error) {
	// TODO: wait to be synced

	topic, ok := c.topics.Load(topicName)
	if !ok {
		return nil, nil
	}

	return topic, nil
}

func (c *Cluster) APIUpdateTopic(topicName string, partitions int32, config map[string]*string) (*topics.Topic, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	topic, ok := c.topics.Load(topicName)
	if !ok {
		return nil, fmt.Errorf("topic does not exist")
	}

	adminClient := kadm.NewClient(c.franzKafkaClient)
	if topic.Partitions > partitions {
		// update partitions
		resp, err := adminClient.UpdatePartitions(ctx, int(partitions), topicName)
		if err != nil {
			return nil, err
		}
		if err := resp[topicName].Err; err != nil {
			return nil, fmt.Errorf("error from kafka while updating partitions: %w", err)
		}
	} else if topic.Partitions < partitions {
		return nil, fmt.Errorf("cannot decrease topic partitions")
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
			return nil, err
		}
		for _, r := range resp {
			if r.Err != nil {
				return nil, fmt.Errorf("error from kafka while updating config: %w", err)
			}
		}
	}

	topic = &topics.Topic{
		Name:       topicName,
		Partitions: partitions,
		Config:     config,
	}

	topicMessage := &TopicMessage{
		Name:    topicName,
		Action:  TopicMessageActionUpdate,
		Cluster: c.Name,
		Topic:   topic,
	}

	kafkaClient := c.controller.franzKafkaClient
	topicMessageBytes, _ := json.Marshal(topicMessage)
	record := kgo.KeySliceRecord([]byte(fmt.Sprintf("cluster-%s-topic-%s", c.Name, topic.Name)), topicMessageBytes)
	record.Topic = InternalTopicTopicConfig
	produceResp := kafkaClient.ProduceSync(context.TODO(), record)
	if produceResp.FirstErr() != nil {
		return nil, produceResp.FirstErr()
	}

	return topic, nil
}

func (c *Cluster) APIDeleteTopicCluster(topicName string) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

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

	topicMessage := &TopicMessage{
		Name:    topicName,
		Action:  TopicMessageActionDelete,
		Cluster: c.Name,
	}

	kafkaClient := c.controller.franzKafkaClient
	topicMessageBytes, _ := json.Marshal(topicMessage)
	record := kgo.KeySliceRecord([]byte(fmt.Sprintf("cluster-%s-topic-%s", c.Name, topicName)), topicMessageBytes)
	record.Topic = InternalTopicTopicConfig
	produceResp := kafkaClient.ProduceSync(context.TODO(), record)
	if produceResp.FirstErr() != nil {
		return produceResp.FirstErr()
	}

	return nil
}
