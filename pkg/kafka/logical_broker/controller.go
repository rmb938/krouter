package logical_broker

import (
	"sync"

	"github.com/go-logr/logr"
	"github.com/puzpuzpuz/xsync"
	"github.com/rmb938/krouter/pkg/redisw"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Controller struct {
	log logr.Logger

	topicPointerSyncOnce sync.Once
	syncedChan           chan struct{}

	redisClient *redisw.RedisClient

	authorizer *Authorizer

	kafkaAddrs       []string
	franzKafkaClient *kgo.Client

	topicPointers *xsync.MapOf[string]
}

func NewController(log logr.Logger, addrs []string, redisClient *redisw.RedisClient) (*Controller, error) {
	controller := &Controller{
		log:        log.WithName("controller"),
		syncedChan: make(chan struct{}),

		kafkaAddrs:    addrs,
		redisClient:   redisClient,
		topicPointers: xsync.NewMapOf[string](),
	}

	err := controller.initFranzKafkaClient()
	if err != nil {
		return nil, err
	}

	authorizerKafkaClient, err := controller.newFranzKafkaClient(InternalTopicAcls)
	if err != nil {
		return nil, err
	}
	controller.authorizer, err = NewAuthorizer(log, authorizerKafkaClient)
	if err != nil {
		return nil, err
	}

	return controller, nil
}

func (c *Controller) newFranzKafkaClient(topics ...string) (*kgo.Client, error) {
	franzKafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(c.kafkaAddrs...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ConsumeTopics(topics...),
	)
	if err != nil {
		return nil, err
	}

	return franzKafkaClient, nil
}

func (c *Controller) initFranzKafkaClient() error {
	var err error
	c.franzKafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(c.kafkaAddrs...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchCompression(kgo.GzipCompression()),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) waitSynced() {
	<-c.syncedChan
}

func (c *Controller) Start() error {
	err := c.CreateInternalTopics()
	if err != nil {
		return err
	}

	err = c.authorizer.CreateTestACLs()
	if err != nil {
		return err
	}

	go c.ConsumeTopicPointers()
	c.waitSynced()
	c.log.Info("Controller Synced")

	go c.authorizer.ConsumeAcls()
	c.authorizer.WaitSynced()
	c.log.Info("Authorizer Synced")
	return nil
}

func (c *Controller) GetAuthorizer() *Authorizer {
	return c.authorizer
}
