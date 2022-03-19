package logical_broker

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/redisw"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

const (
	ClusterTopicLeaderRedisKeyFmt = "cluster-%s-{topic-%s}-partition-%d-leader"
)

type Cluster struct {
	Name string
	log  logr.Logger

	metadatRefreshCtxCancel context.CancelFunc
	metadataRefreshCtx      context.Context

	redisClient *redisw.RedisClient

	saramaKafkaClient sarama.Client
	franzKafkaClient  *kgo.Client

	topics map[string]*topics.Topic
}

func NewCluster(name string, addrs []string, log logr.Logger, redisClient *redisw.RedisClient) (*Cluster, error) {
	log = log.WithName(fmt.Sprintf("cluster-%s", name))

	metadatRefreshCtx, metadatRefreshCtxCancel := context.WithCancel(context.Background())

	cluster := &Cluster{
		Name:                    name,
		log:                     log,
		redisClient:             redisClient,
		metadatRefreshCtxCancel: metadatRefreshCtxCancel,
		metadataRefreshCtx:      metadatRefreshCtx,
		topics:                  map[string]*topics.Topic{},
	}

	if err := cluster.initSaramaKafkaClient(addrs); err != nil {
		return nil, err
	}

	if err := cluster.initFranzKafkaClient(addrs); err != nil {
		return nil, err
	}

	return cluster, nil
}

func (c *Cluster) initSaramaKafkaClient(addrs []string) error {
	c.log.Info("Creating cluster client")
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_6_0_0

	// TODO: if the brokers we are connected to restart we get a broken pipe error
	// {"level":"error","ts":1646785492.4811382,"logger":"router.packet-processor.describe-groups-v0-handler","msg":"Error describing group to controller","from-address":"127.0.0.1:43656","group":"same-group1","error":"write tcp [::1]:45220->[::1]:19093: write: broken pipe"}
	// {"level":"error","ts":1646785492.481163,"logger":"router","msg":"error processing packet","from-address":"127.0.0.1:43656","error":"error handling packet: error describing group to controller: write tcp [::1]:45220->[::1]:19093: write: broken pipe"}
	// TODO: we need to come up with a way to check the connection and re-open it if needed
	//  it eventually fixes itself but we can't tolerate these

	var err error
	c.saramaKafkaClient, err = sarama.NewClient(addrs, saramaConfig)
	if err != nil {
		return fmt.Errorf("error creating kafka client for cluster: %v: %w", c.Name, err)
	}

	if err := c.saramaKafkaClient.RefreshMetadata(); err != nil {
		return fmt.Errorf("error refreshing metadata for cluster: %v: %w", c.Name, err)
	}

	if _, err := c.saramaKafkaClient.RefreshController(); err != nil {
		return fmt.Errorf("error refreshing controller for cluster: %v: %w", c.Name, err)
	}

	return nil
}

func (c *Cluster) initFranzKafkaClient(addrs []string) error {

	var err error
	c.franzKafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(addrs...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) Close() error {
	c.metadatRefreshCtxCancel()
	c.franzKafkaClient.Close()
	return c.saramaKafkaClient.Close()
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

func (c *Cluster) TopicLeader(ctx context.Context, topic string, partition int32) (int, error) {
	brokerID, err := c.redisClient.Client.Get(ctx, fmt.Sprintf(ClusterTopicLeaderRedisKeyFmt, c.Name, topic, partition)).Int()
	if err != nil {
		if err == redis.Nil {
			return -1, nil
		}
		return -1, err
	}

	return brokerID, nil
}

func (c *Cluster) TopicMetadata(ctx context.Context, topics []string) (*kmsg.MetadataResponse, error) {

	metadataRequest := kmsg.NewPtrMetadataRequest()
	metadataRequest.AllowAutoTopicCreation = true
	for _, topic := range topics {
		metadataTopicRequest := kmsg.NewMetadataRequestTopic()
		metadataTopicRequest.Topic = &topic
		metadataRequest.Topics = append(metadataRequest.Topics, metadataTopicRequest)
	}

	resp, err := c.franzKafkaClient.Request(ctx, metadataRequest)
	if err != nil {
		return nil, err
	}

	metadataResponse := resp.(*kmsg.MetadataResponse)

	for _, topic := range metadataResponse.Topics {
		if topic.ErrorCode != int16(errors.None) {
			continue
		}

		for _, partition := range topic.Partitions {
			if partition.ErrorCode != int16(errors.None) {
				continue
			}

			// Clients typically refresh metadata every 10 minutes
			// so expiring in an hour "should" be good enough
			err := c.redisClient.Client.Set(ctx, fmt.Sprintf(ClusterTopicLeaderRedisKeyFmt, c.Name, *topic.Topic, partition.Partition), partition.Leader, 1*time.Hour).Err()
			if err != nil {
				return nil, err
			}
		}
	}

	return metadataResponse, nil
}

func (c *Cluster) FranzProduce(topic *topics.Topic, partition int32, transactionID *string, timeoutMillis int32, recordBytes []byte) (*kmsg.ProduceResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	brokerID, err := c.TopicLeader(ctx, topic.Name, partition)
	if err != nil {
		return nil, err
	}

	if brokerID == -1 {
		// can't find topic leader
		return &kmsg.ProduceResponse{
			Topics: []kmsg.ProduceResponseTopic{
				{
					Topic: topic.Name,
					Partitions: []kmsg.ProduceResponseTopicPartition{
						{
							Partition: partition,
							ErrorCode: int16(errors.NotLeaderOrFollower),
						},
					},
				},
			},
		}, nil
	}

	response, err := c.franzKafkaClient.Broker(brokerID).RetriableRequest(ctx, &kmsg.ProduceRequest{
		TransactionID: transactionID,
		TimeoutMillis: timeoutMillis,
		Topics: []kmsg.ProduceRequestTopic{{
			Topic: topic.Name,
			Partitions: []kmsg.ProduceRequestTopicPartition{
				{
					Partition: partition,
					Records:   recordBytes,
				},
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.ProduceResponse), nil
}

func (c *Cluster) SaramaProduce(topic *topics.Topic, partition int32, request *sarama.ProduceRequest) (*sarama.ProduceResponse, error) {
	leaderID, err := c.TopicLeader(context.TODO(), topic.Name, partition)
	if err != nil {
		return nil, err
	}

	broker, err := c.saramaKafkaClient.Broker(int32(leaderID))
	if err != nil {
		return nil, err
	}
	return broker.Produce(request)
}

func (c *Cluster) ListOffsets(topic *topics.Topic, partition int32, request *sarama.OffsetRequest) (*sarama.OffsetResponse, error) {
	leaderID, err := c.TopicLeader(context.TODO(), topic.Name, partition)
	if err != nil {
		return nil, err
	}

	broker, err := c.saramaKafkaClient.Broker(int32(leaderID))
	if err != nil {
		return nil, err
	}

	return broker.GetAvailableOffsets(request)
}

func (c *Cluster) Fetch(topic *topics.Topic, partition int32, request *sarama.FetchRequest) (*sarama.FetchResponse, error) {
	leaderID, err := c.TopicLeader(context.TODO(), topic.Name, partition)
	if err != nil {
		return nil, err
	}

	broker, err := c.saramaKafkaClient.Broker(int32(leaderID))
	if err != nil {
		return nil, err
	}

	return broker.Fetch(request)
}
