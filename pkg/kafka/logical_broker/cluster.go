package logical_broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	"github.com/outcaste-io/ristretto"
	"github.com/puzpuzpuz/xsync"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

const (
	ClusterTopicLeaderKeyFmt = "topic-%s-partition-%d"
)

type Cluster struct {
	Name string
	log  logr.Logger

	topicConfigSyncedOnce sync.Once
	topicLeaderSyncedOnce sync.Once
	syncedChan            chan struct{}

	controller *Controller

	franzKafkaClient *kgo.Client

	topics           *xsync.MapOf[*topics.Topic]
	topicLeaderCache *ristretto.Cache
}

func NewCluster(name string, addrs []string, log logr.Logger, controller *Controller) (*Cluster, error) {
	log = log.WithName(fmt.Sprintf("cluster-%s", name))

	cluster := &Cluster{
		Name:       name,
		log:        log,
		syncedChan: make(chan struct{}),
		controller: controller,
		topics:     xsync.NewMapOf[*topics.Topic](),
	}

	if err := cluster.initFranzKafkaClient(addrs); err != nil {
		return nil, err
	}

	var err error
	cluster.topicLeaderCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (c *Cluster) initFranzKafkaClient(addrs []string) error {
	var err error
	c.franzKafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(addrs...),
		kgo.RequiredAcks(kgo.AllISRAcks()), // Required to support acks all, so let's just upgrade everyone
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) WaitSynced() {
	// Two loads, one for topic configs and one for leaders
	<-c.syncedChan
	<-c.syncedChan
}

func (c *Cluster) Close() error {
	c.franzKafkaClient.Close()
	return nil
}

func (c *Cluster) topicLeader(topic string, partition int32) (int, error) {
	brokerIDInterf, found := c.topicLeaderCache.Get(fmt.Sprintf(ClusterTopicLeaderKeyFmt, topic, partition))
	if !found {
		return -1, nil
	}
	return brokerIDInterf.(int), nil
}

func (c *Cluster) TopicMetadata(ctx context.Context, topics []string) (*kmsg.MetadataResponse, error) {
	kafkaClient := c.controller.franzKafkaClient

	metadataRequest := kmsg.NewPtrMetadataRequest()
	metadataRequest.AllowAutoTopicCreation = false // don't allow topic creation
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

			// publish topic partition leader
			topicLeaderMessage := &TopicLeaderMessage{
				Name:      *topic.Topic,
				Cluster:   c.Name,
				Partition: partition.Partition,
				Leader:    int(partition.Leader),
			}
			topicLeaderMessageBytes, _ := json.Marshal(topicLeaderMessage)
			record := kgo.KeySliceRecord([]byte(fmt.Sprintf("cluster-%s-topic-%s-partition-%d", c.Name, *topic.Topic, partition.Partition)), topicLeaderMessageBytes)
			record.Topic = InternalTopicTopicLeader
			produceResp := kafkaClient.ProduceSync(ctx, record)
			if produceResp.FirstErr() != nil {
				return nil, produceResp.FirstErr()
			}
		}
	}

	return metadataResponse, nil
}

func (c *Cluster) Produce(transactionID *string, timeoutMillis int32, topics []kmsg.ProduceRequestTopic) (*kmsg.ProduceResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// get leader from first topic since all topics and partitions in this request should be to the same leader
	//  if our cache has the wrong one some (or all) of our topics will error and the client will refresh metadata
	leaderID, err := c.topicLeader(topics[0].Topic, topics[0].Partitions[0].Partition)
	if err != nil {
		return nil, err
	}

	if leaderID == -1 {
		resp := kmsg.NewPtrProduceResponse()

		for _, topic := range topics {
			respTopic := kmsg.NewProduceResponseTopic()
			respTopic.Topic = topic.Topic

			for _, partition := range topic.Partitions {
				respPartition := kmsg.NewProduceResponseTopicPartition()
				respPartition.Partition = partition.Partition
				respPartition.ErrorCode = int16(errors.NotLeaderOrFollower)

				respTopic.Partitions = append(respTopic.Partitions, respPartition)
			}

			resp.Topics = append(resp.Topics, respTopic)
		}

		return resp, err
	}

	req := &kmsg.ProduceRequest{
		TransactionID: transactionID,
		TimeoutMillis: timeoutMillis,
		Topics:        topics,
	}

	response, err := c.franzKafkaClient.Broker(leaderID).RetriableRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.ProduceResponse), nil
}

func (c *Cluster) ListOffsets(topicName string, partition int32, request *kmsg.ListOffsetsRequest) (*kmsg.ListOffsetsResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	leaderID, err := c.topicLeader(topicName, partition)
	if err != nil {
		return nil, err
	}

	if leaderID == -1 {
		// can't find topic partition leader
		response := kmsg.NewPtrListOffsetsResponse()

		responseTopic := kmsg.NewListOffsetsResponseTopic()
		responseTopic.Topic = topicName

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

func (c *Cluster) Fetch(topicName string, partition int32, request *kmsg.FetchRequest) (*kmsg.FetchResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	leaderID, err := c.topicLeader(topicName, partition)
	if err != nil {
		return nil, err
	}

	if leaderID == -1 {
		// can't find topic leader
		response := kmsg.NewPtrFetchResponse()
		response.SessionID = request.SessionID

		responseTopic := kmsg.NewFetchResponseTopic()
		responseTopic.Topic = topicName

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
