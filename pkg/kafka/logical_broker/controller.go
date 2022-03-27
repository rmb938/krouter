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

	syncedOnce sync.Once
	syncedChan chan struct{}

	redisClient *redisw.RedisClient

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

func (c *Controller) WaitSynced() {
	<-c.syncedChan
}
