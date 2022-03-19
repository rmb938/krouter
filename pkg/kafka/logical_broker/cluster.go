package logical_broker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
	"github.com/rmb938/krouter/pkg/redisw"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/kversion"
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

	franzKafkaClient *kgo.Client

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

	if err := cluster.initFranzKafkaClient(addrs); err != nil {
		return nil, err
	}

	return cluster, nil
}

func (c *Cluster) initFranzKafkaClient(addrs []string) error {

	// pull maxVersions
	maxVersions := kversion.Stable()
	// need to pin the max version of sync_group due to protocol setting in newer versions
	maxVersions.SetMaxKeyVersion(sync_group.Key, implSyncGroupV3.Version)

	var err error
	c.franzKafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(addrs...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.MaxVersions(maxVersions),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) Close() error {
	c.metadatRefreshCtxCancel()
	c.franzKafkaClient.Close()
	return nil
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

func (c *Cluster) topicLeader(ctx context.Context, topic string, partition int32) (int, error) {
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
	metadataRequest.Topics = make([]kmsg.MetadataRequestTopic, 0)
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

func (c *Cluster) Produce(topic *topics.Topic, partition int32, transactionID *string, timeoutMillis int32, recordBytes []byte) (*kmsg.ProduceResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	leaderID, err := c.topicLeader(ctx, topic.Name, partition)
	if err != nil {
		return nil, err
	}

	if leaderID == -1 {
		// can't find topic partition leader
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

	response, err := c.franzKafkaClient.Broker(leaderID).RetriableRequest(ctx, &kmsg.ProduceRequest{
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

func (c *Cluster) ListOffsets(topic *topics.Topic, partition int32, request *kmsg.ListOffsetsRequest) (*kmsg.ListOffsetsResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	leaderID, err := c.topicLeader(ctx, topic.Name, partition)
	if err != nil {
		return nil, err
	}

	if leaderID == -1 {
		// can't find topic partition leader
		response := kmsg.NewPtrListOffsetsResponse()

		responseTopic := kmsg.NewListOffsetsResponseTopic()
		responseTopic.Topic = topic.Name

		responsePartition := kmsg.NewListOffsetsResponseTopicPartition()
		responsePartition.Partition = partition
		responsePartition.ErrorCode = int16(errors.UnknownTopicOrPartition)

		responseTopic.Partitions = append(responseTopic.Partitions, responsePartition)

		response.Topics = append(response.Topics, responseTopic)
		return response, nil
	}

	response, err := c.franzKafkaClient.Broker(leaderID).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.ListOffsetsResponse), nil
}

func (c *Cluster) Fetch(topic *topics.Topic, partition int32, request *kmsg.FetchRequest) (*kmsg.FetchResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	leaderID, err := c.topicLeader(ctx, topic.Name, partition)
	if err != nil {
		return nil, err
	}

	if leaderID == -1 {
		// can't find topic leader
		response := kmsg.NewPtrFetchResponse()
		response.SessionID = request.SessionID

		responseTopic := kmsg.NewFetchResponseTopic()
		responseTopic.Topic = topic.Name

		responseTopicPartition := kmsg.NewFetchResponseTopicPartition()
		responseTopicPartition.Partition = partition
		responseTopicPartition.ErrorCode = int16(errors.UnknownTopicOrPartition)

		responseTopic.Partitions = append(responseTopic.Partitions, responseTopicPartition)

		response.Topics = append(response.Topics, responseTopic)

		return response, nil
	}

	response, err := c.franzKafkaClient.Broker(leaderID).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.FetchResponse), nil
}
