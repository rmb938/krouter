package logical_broker

import (
	"fmt"
	"net"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker/topics"
	"github.com/rmb938/krouter/pkg/redisw"
)

type Broker struct {
	AdvertiseListener *net.TCPAddr
	ClusterID         string

	log logr.Logger

	ephemeralID string

	redisClient *redisw.RedisClient

	controller *Controller
	clusters   map[string]*Cluster
}

func InitBroker(log logr.Logger, advertiseListener *net.TCPAddr, clusterID string, redisAddresses []string) (*Broker, error) {
	log = log.WithName("broker")

	redisClient, err := redisw.NewRedisClient(redisAddresses)
	if err != nil {
		return nil, err
	}

	broker := &Broker{
		AdvertiseListener: advertiseListener,
		ClusterID:         clusterID,

		log:         log,
		ephemeralID: uuid.Must(uuid.NewRandom()).String(),
		redisClient: redisClient,
		clusters:    map[string]*Cluster{},
	}

	return broker, nil
}

func (b *Broker) Close() error {
	var result error

	for _, cluster := range b.clusters {
		err := cluster.Close()
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

func (b *Broker) registerCluster(name string, addrs []string) (*Cluster, error) {
	b.log.Info("Registering cluster", "cluster", name)
	if _, ok := b.clusters[name]; ok {
		return nil, fmt.Errorf("cluster %s is already registered", name)
	}

	var err error
	b.clusters[name], err = NewCluster(name, addrs, b.log, b.redisClient)
	if err != nil {
		return nil, err
	}

	return b.clusters[name], nil
}

func (b *Broker) InitClusters() error {
	// TODO: load from config or env vars

	cluster, err := b.registerCluster("controller", []string{"localhost:19093"})
	if err != nil {
		return err
	}

	b.controller, err = NewController(b.log, cluster)
	if err != nil {
		return err
	}

	cluster, err = b.registerCluster("cluster1", []string{"localhost:9093"})
	if err != nil {
		return err
	}

	// TODO: this is temporary
	cluster.AddTopic(&topics.Topic{
		Name:       "test1",
		Partitions: 1,
		Config:     nil,
		Enabled:    true,
	})

	cluster, err = b.registerCluster("cluster2", []string{"localhost:9094"})
	if err != nil {
		return err
	}

	// TODO: this is temporary
	cluster.AddTopic(&topics.Topic{
		Name:       "test2",
		Partitions: 1,
		Config:     nil,
		Enabled:    true,
	})

	cluster, err = b.registerCluster("cluster3", []string{"localhost:9392"})
	if err != nil {
		return err
	}

	// TODO: this is temporary
	cluster.AddTopic(&topics.Topic{
		Name:       "test3",
		Partitions: 1,
		Config:     nil,
		Enabled:    true,
	})

	return nil
}

func (b *Broker) GetController() *Controller {
	return b.controller
}

func (b *Broker) GetClusters() map[string]*Cluster {
	return b.clusters
}

func (b *Broker) GetTopics() []*topics.Topic {
	var allTopics []*topics.Topic

	for _, cluster := range b.clusters {
		for _, value := range cluster.GetTopics() {
			allTopics = append(allTopics, value)
		}
	}

	return allTopics
}

func (b *Broker) GetTopic(name string) (*Cluster, *topics.Topic) {
	for _, cluster := range b.clusters {
		if topic := cluster.GetTopic(name); topic != nil && topic.Enabled {
			return cluster, topic
		}
	}

	return nil, nil
}
