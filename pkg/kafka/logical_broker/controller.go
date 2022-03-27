package logical_broker

import (
	"sync"

	"github.com/go-logr/logr"
	"github.com/puzpuzpuz/xsync"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	"github.com/rmb938/krouter/pkg/redisw"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
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
	// pull maxVersions
	maxVersions := kversion.Stable()
	// need to pin the max version of sync_group due to protocol setting in newer versions
	// we support version 3, version 4 just adds tagged fields
	// TODO: figure out how to not need this
	maxVersions.SetMaxKeyVersion(sync_group.Key, 4)

	var err error
	c.franzKafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(c.kafkaAddrs...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchCompression(kgo.GzipCompression()),
		kgo.MaxVersions(maxVersions),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) WaitSynced() {
	<-c.syncedChan
}
