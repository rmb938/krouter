package logical_broker

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
)

type Cluster struct {
	Name        string
	log         logr.Logger
	kafkaClient sarama.Client

	topics map[string]*topics.Topic
}

func NewCluster(name string, addrs []string, log logr.Logger) (*Cluster, error) {
	log = log.WithName(fmt.Sprintf("cluster-%s", name))

	cluster := &Cluster{
		Name:   name,
		log:    log,
		topics: map[string]*topics.Topic{},
	}

	log.Info("Creating cluster client")
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_6_0_0

	var err error
	cluster.kafkaClient, err = sarama.NewClient(addrs, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka client for cluster: %v: %w", name, err)
	}

	if err := cluster.kafkaClient.RefreshMetadata(); err != nil {
		return nil, fmt.Errorf("error refreshing metadata for cluster: %v: %w", name, err)
	}

	if _, err := cluster.kafkaClient.RefreshController(); err != nil {
		return nil, fmt.Errorf("error refreshing controller for cluster: %v: %w", name, err)
	}

	return cluster, nil
}

func (c *Cluster) Close() error {
	return c.kafkaClient.Close()
}

func (c *Cluster) GetTopics() map[string]*topics.Topic {
	topics := map[string]*topics.Topic{}

	for name, value := range c.topics {
		topics[name] = value.Clone()
	}

	return topics
}

func (c *Cluster) GetTopic(name string) *topics.Topic {
	if topic, ok := c.topics[name]; ok {
		return topic.Clone()
	}

	return nil
}

func (c *Cluster) AddTopic(topic *topics.Topic) {
	c.topics[topic.Name] = topic.Clone()
}

func (c *Cluster) RemoveTopic(topic *topics.Topic) {
	delete(c.topics, topic.Name)
}

func (c *Cluster) Produce(topic *topics.Topic, partition int32, request *sarama.ProduceRequest) (*sarama.ProduceResponse, error) {
	broker, err := c.kafkaClient.Leader(topic.Name, partition)
	if err != nil {
		return nil, err
	}

	return broker.Produce(request)
}
