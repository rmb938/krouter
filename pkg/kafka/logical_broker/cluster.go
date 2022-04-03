package logical_broker

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	"github.com/puzpuzpuz/xsync"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type Cluster struct {
	Name string
	log  logr.Logger

	topicConfigSyncedOnce sync.Once
	syncedChan            chan struct{}

	controller *Controller

	franzKafkaClient *kgo.Client

	topics *xsync.MapOf[*models.Topic]
}

func NewCluster(name string, addrs []string, log logr.Logger, controller *Controller) (*Cluster, error) {
	log = log.WithName(fmt.Sprintf("cluster-%s", name))

	cluster := &Cluster{
		Name:       name,
		log:        log,
		syncedChan: make(chan struct{}),
		controller: controller,
		topics:     xsync.NewMapOf[*models.Topic](),
	}

	if err := cluster.initFranzKafkaClient(addrs); err != nil {
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
}

func (c *Cluster) Close() error {
	c.franzKafkaClient.Close()
	return nil
}

func (c *Cluster) TopicMetadata(ctx context.Context, topics []string) (*kmsg.MetadataResponse, error) {
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

	return resp.(*kmsg.MetadataResponse), nil
}

func (c *Cluster) Produce(brokerID int32, transactionID *string, timeoutMillis int32, topics []kmsg.ProduceRequestTopic) (*kmsg.ProduceResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	req := &kmsg.ProduceRequest{
		TransactionID: transactionID,
		TimeoutMillis: timeoutMillis,
		Topics:        topics,
	}

	response, err := c.franzKafkaClient.Broker(int(brokerID)).RetriableRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.ProduceResponse), nil
}

func (c *Cluster) ListOffsets(brokerID int32, request *kmsg.ListOffsetsRequest) (*kmsg.ListOffsetsResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	response, err := c.franzKafkaClient.Broker(int(brokerID)).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.ListOffsetsResponse), nil
}

func (c *Cluster) Fetch(brokerID int32, request *kmsg.FetchRequest) (*kmsg.FetchResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	response, err := c.franzKafkaClient.Broker(int(brokerID)).RetriableRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.FetchResponse), nil
}
