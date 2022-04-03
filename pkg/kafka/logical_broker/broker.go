package logical_broker

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/puzpuzpuz/xsync"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/internal_topics_pb"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/models"
	"github.com/rmb938/krouter/pkg/redisw"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

type LogicalBroker struct {
	AdvertiseListener *net.TCPAddr
	ClusterID         string
	BrokerID          int32

	log logr.Logger

	ephemeralID string

	redisClient *redisw.RedisClient

	controller *Controller
	clusters   map[string]*Cluster

	brokerSyncedOnce sync.Once
	syncedChan       chan struct{}

	brokers *xsync.MapOf[*models.Broker]
}

func InitBroker(log logr.Logger, advertiseListener *net.TCPAddr, clusterID string, brokerID int32, redisAddresses []string) (*LogicalBroker, error) {
	log = log.WithName("broker")

	redisClient, err := redisw.NewRedisClient(redisAddresses)
	if err != nil {
		return nil, err
	}

	broker := &LogicalBroker{
		AdvertiseListener: advertiseListener,
		ClusterID:         clusterID,
		BrokerID:          brokerID,

		log:         log,
		syncedChan:  make(chan struct{}),
		ephemeralID: uuid.Must(uuid.NewRandom()).String(),
		redisClient: redisClient,
		clusters:    map[string]*Cluster{},
		brokers:     xsync.NewMapOf[*models.Broker](),
	}

	return broker, nil
}

func (b *LogicalBroker) Close() error {
	var result error

	for _, cluster := range b.clusters {
		err := cluster.Close()
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

func (b *LogicalBroker) GetBroker(brokerID int32) *models.Broker {
	broker, ok := b.brokers.Load(strconv.FormatInt(int64(brokerID), 10))
	if !ok {
		return nil
	}

	return broker
}

func (b *LogicalBroker) GetRegisteredBrokers() []*models.Broker {
	var brokers []*models.Broker

	b.brokers.Range(func(_ string, broker *models.Broker) bool {
		brokers = append(brokers, broker)
		return true
	})

	return brokers
}

func (b *LogicalBroker) registerCluster(name string, addrs []string) (*Cluster, error) {
	b.log.Info("Registering cluster", "cluster", name)
	if _, ok := b.clusters[name]; ok {
		return nil, fmt.Errorf("cluster %s is already registered", name)
	}

	var err error
	b.clusters[name], err = NewCluster(name, addrs, b.log, b.controller)
	if err != nil {
		return nil, err
	}

	return b.clusters[name], nil
}

func (b *LogicalBroker) waitSynced() {
	<-b.syncedChan
}

func (b *LogicalBroker) produceHeartbeat() error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	brokerKey := &internal_topics_pb.BrokerMessageKey{
		BrokerId: b.BrokerID,
	}

	brokerValue := &internal_topics_pb.BrokerMessageValue{
		Endpoint: &internal_topics_pb.BrokerEndpoint{
			Host: b.AdvertiseListener.IP.String(),
			Port: int32(b.AdvertiseListener.Port),
		},
		Rack: "rack",
	}

	brokerMessageKeyBytes, err := proto.Marshal(brokerKey)
	if err != nil {
		return err
	}

	brokerMessageValueBytes, err := proto.Marshal(brokerValue)
	if err != nil {
		return err
	}

	record := kgo.KeySliceRecord(brokerMessageKeyBytes, brokerMessageValueBytes)
	record.Topic = InternalTopicBrokers
	resp := b.controller.franzKafkaClient.ProduceSync(ctx, record)
	if resp.FirstErr() != nil {
		return resp.FirstErr()
	}

	return nil
}

func (b *LogicalBroker) Start() error {
	if err := b.initClusters(); err != nil {
		return err
	}

	// produce heartbeat every 5 seconds
	go func() {
		// TODO: wait until healthy to start ticking
		//  if we aren't healthy and we tick clients won't be able to connect
		ticker := time.NewTicker(5 * time.Second)
		for ; true; <-ticker.C {
			err := b.produceHeartbeat()
			if err != nil {
				b.log.Error(err, "error producing broker heartbeat")
			}
		}
	}()

	go b.ConsumeBrokers()

	// loop over hashring every 10 seconds and remove brokers that haven't heartbeat in 30 seconds
	go func() {
		ticker := time.NewTicker(10 * time.Second)

		for ; true; <-ticker.C {
			for _, broker := range b.GetRegisteredBrokers() {
				if time.Now().Sub(broker.LastHeartbeat) > 30*time.Second {
					b.brokers.Delete(strconv.FormatInt(int64(broker.ID), 10))
				}
			}
		}
	}()

	b.waitSynced()
	b.log.Info("Logical Broker Synced")

	return nil
}

func (b *LogicalBroker) initClusters() error {
	// TODO: load from config or env vars

	var err error
	b.controller, err = NewController(b.log, []string{"localhost:19093"}, b.redisClient)
	if err != nil {
		return err
	}

	err = b.controller.Start()
	if err != nil {
		return err
	}
	b.log.Info("Controller Started")

	// TODO: remove this it's temporary
	pointer, err := b.controller.APIGetTopicPointer("test1")
	if err != nil {
		return err
	}
	if pointer == nil {
		err = b.controller.APISetTopicPointer("test1", "cluster1")
		if err != nil {
			return err
		}
	}

	// TODO: remove this it's temporary
	pointer, err = b.controller.APIGetTopicPointer("test2")
	if err != nil {
		return err
	}
	if pointer == nil {
		err = b.controller.APISetTopicPointer("test2", "cluster2")
		if err != nil {
			return err
		}
	}

	// TODO: remove this it's temporary
	pointer, err = b.controller.APIGetTopicPointer("test3")
	if err != nil {
		return err
	}
	if pointer == nil {
		err = b.controller.APISetTopicPointer("test3", "cluster3")
		if err != nil {
			return err
		}
	}

	cluster1, err := b.registerCluster("cluster1", []string{"localhost:9093"})
	if err != nil {
		return err
	}

	cluster2, err := b.registerCluster("cluster2", []string{"localhost:9094"})
	if err != nil {
		return err
	}

	cluster3, err := b.registerCluster("cluster3", []string{"localhost:9392"})
	if err != nil {
		return err
	}

	for _, cluster := range b.clusters {
		go cluster.ConsumeTopicConfigs()
		cluster.WaitSynced()
		b.log.Info("Cluster Synced", "cluster", cluster.Name)
	}

	// TODO: remove this it's temporary
	topic, err := cluster1.APIGetTopic("test1")
	if err != nil {
		return err
	}
	if topic == nil {
		_, err := cluster1.APICreateTopic("test1", 1, 1, map[string]*string{})
		if err != nil {
			return fmt.Errorf("err creating topic test1: %w", err)
		}
	}

	// TODO: remove this it's temporary
	topic, err = cluster2.APIGetTopic("test2")
	if err != nil {
		return err
	}
	if topic == nil {
		_, err := cluster2.APICreateTopic("test2", 1, 1, map[string]*string{})
		if err != nil {
			return fmt.Errorf("err creating topic test2: %w", err)
		}
	}

	// TODO: remove this it's temporary
	topic, err = cluster3.APIGetTopic("test3")
	if err != nil {
		return err
	}
	if topic == nil {
		_, err := cluster3.APICreateTopic("test3", 3, 3, map[string]*string{})
		if err != nil {
			return fmt.Errorf("err creating topic test3: %w", err)
		}
	}

	return nil
}

func (b *LogicalBroker) GetController() *Controller {
	return b.controller
}

func (b *LogicalBroker) GetClusters() map[string]*Cluster {
	return b.clusters
}

func (b *LogicalBroker) GetTopics() ([]*models.Topic, error) {
	var allTopics []*models.Topic

	for _, cluster := range b.clusters {
		cluster.topics.Range(func(topicName string, topic *models.Topic) bool {
			if clusterName, ok := b.controller.topicPointers.Load(topicName); ok {
				if clusterName == cluster.Name {
					allTopics = append(allTopics, topic)
				}
			}
			return true
		})
	}

	return allTopics, nil
}

func (b *LogicalBroker) GetClusterByTopic(topicName string) *Cluster {
	if clusterName, ok := b.controller.topicPointers.Load(topicName); ok {
		if cluster, ok := b.clusters[clusterName]; ok {
			return cluster
		}
	}

	return nil
}

func (b *LogicalBroker) GetTopic(topicName string) (*Cluster, *models.Topic) {
	if cluster := b.GetClusterByTopic(topicName); cluster != nil {
		if topic, ok := cluster.topics.Load(topicName); ok {
			return cluster, topic
		}
	}

	return nil, nil
}
