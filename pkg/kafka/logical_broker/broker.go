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
	b.clusters[name], err = NewCluster(name, addrs, b.log, b.controller)
	if err != nil {
		return nil, err
	}

	return b.clusters[name], nil
}

func (b *Broker) InitClusters() error {
	// TODO: load from config or env vars

	var err error
	b.controller, err = NewController(b.log, []string{"localhost:19093"}, b.redisClient)
	if err != nil {
		return err
	}

	err = b.controller.CreateInternalTopics()
	if err != nil {
		return err
	}

	go b.controller.ConsumeTopicPointers()
	b.controller.WaitSynced()
	b.log.Info("Controller Synced")

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
		go cluster.ConsumeTopicLeaders()
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
		_, err := cluster3.APICreateTopic("test3", 1, 3, map[string]*string{})
		if err != nil {
			return fmt.Errorf("err creating topic test3: %w", err)
		}
	}

	return nil
}

func (b *Broker) GetController() *Controller {
	return b.controller
}

func (b *Broker) GetClusters() map[string]*Cluster {
	return b.clusters
}

func (b *Broker) GetTopics() ([]*topics.Topic, error) {
	var allTopics []*topics.Topic

	for _, cluster := range b.clusters {
		cluster.topics.Range(func(topicName string, topic *topics.Topic) bool {
			if clusterName, ok := b.controller.topicPointers.Load(topicName); ok {
				if clusterName == cluster.Name {
					allTopics = append(allTopics, topic)
				}
			}
			return false
		})
	}

	return allTopics, nil
}

func (b *Broker) GetTopic(topicName string) (*Cluster, *topics.Topic) {

	if clusterName, ok := b.controller.topicPointers.Load(topicName); ok {
		if cluster, ok := b.clusters[clusterName]; ok {
			topic, ok := cluster.topics.Load(topicName)
			if ok {
				return cluster, topic
			}
		}
	}

	return nil, nil
}
