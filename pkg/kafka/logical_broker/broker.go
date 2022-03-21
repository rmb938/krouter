package logical_broker

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
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

	cluster, err := b.registerCluster("cluster1", []string{"localhost:9093"})
	if err != nil {
		return err
	}
	// TODO: remove this it's temporary
	topic, err := cluster.APIGetTopic("test1")
	if err != nil {
		return err
	}
	if topic == nil {
		b.log.Info("Creating topic")
		_, err := cluster.APICreateTopic("test1", 1, 1, map[string]*string{})
		if err != nil {
			return err
		}
	}

	cluster, err = b.registerCluster("cluster2", []string{"localhost:9094"})
	if err != nil {
		return err
	}
	// TODO: remove this it's temporary
	topic, err = cluster.APIGetTopic("test2")
	if err != nil {
		return err
	}
	if topic == nil {
		_, err := cluster.APICreateTopic("test2", 1, 1, map[string]*string{})
		if err != nil {
			return err
		}
	}

	cluster, err = b.registerCluster("cluster3", []string{"localhost:9392"})
	if err != nil {
		return err
	}
	// TODO: remove this it's temporary
	topic, err = cluster.APIGetTopic("test3")
	if err != nil {
		return err
	}
	if topic == nil {
		_, err := cluster.APICreateTopic("test3", 1, 3, map[string]*string{})
		if err != nil {
			return err
		}
	}

	cluster, err = b.registerCluster("controller", []string{"localhost:19093"})
	if err != nil {
		return err
	}

	b.controller, err = NewController(b.log, cluster)
	if err != nil {
		return err
	}

	// TODO: remove this it's temporary
	pointer, err := b.controller.APIGetTopicPointer("test1")
	if err != nil {
		return err
	}
	if pointer == nil {
		err = b.controller.APISetTopicPointer("test1", b.clusters["cluster1"])
		if err != nil {
			return err
		}
	}
	pointer, err = b.controller.APIGetTopicPointer("test2")
	if err != nil {
		return err
	}
	if pointer == nil {
		err = b.controller.APISetTopicPointer("test2", b.clusters["cluster2"])
		if err != nil {
			return err
		}
	}
	pointer, err = b.controller.APIGetTopicPointer("test3")
	if err != nil {
		return err
	}
	if pointer == nil {
		err = b.controller.APISetTopicPointer("test3", b.clusters["cluster3"])
		if err != nil {
			return err
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

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer redisContextCancel()

	err := b.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		scan := tx.Scan(redisContext, 0, fmt.Sprintf("%s-*", TopicConfigRedisKeyFmtPrefix), 0).Iterator()

		for scan.Next(redisContext) {
			key := scan.Val()
			pointer, err := tx.Get(redisContext, key).Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				return err
			}

			topicData, err := tx.HGetAll(redisContext, pointer).Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				return err
			}

			partitions, _ := strconv.Atoi(topicData["partitions"])
			config := make(map[string]*string)

			for key, value := range topicData {
				if strings.HasPrefix(key, "config.") {
					configKey := strings.TrimPrefix(key, "config.")
					config[configKey] = &value
				}
			}

			allTopics = append(allTopics, &topics.Topic{
				Name:       strings.TrimPrefix(key, fmt.Sprintf("%s-", TopicConfigRedisKeyFmtPrefix)),
				Partitions: int32(partitions),
				Config:     config,
			})
		}
		return nil
	})

	return allTopics, err
}

func (b *Broker) GetTopic(name string) (*Cluster, *topics.Topic, error) {
	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer redisContextCancel()

	var cluster *Cluster
	var topic *topics.Topic
	topicRedisKey := fmt.Sprintf(TopicConfigRedisKeyFmt, name)
	err := b.redisClient.Client.Watch(redisContext, func(tx *redis.Tx) error {
		pointer, err := tx.Get(redisContext, topicRedisKey).Result()
		if err == redis.Nil {
			return nil
		}
		if err != nil {
			return err
		}

		clusterName := strings.TrimSuffix(strings.TrimPrefix(pointer, fmt.Sprintf("%s-", TopicConfigClusterRedisKeyFmtPrefix)), fmt.Sprintf("-%s", fmt.Sprintf(TopicConfigClusterRedisKeyFmtSuffix, name)))
		cluster, _ = b.clusters[clusterName]

		topicData, err := tx.HGetAll(redisContext, pointer).Result()
		if err == redis.Nil {
			return nil
		}
		if err != nil {
			return err
		}

		partitions, _ := strconv.Atoi(topicData["partitions"])
		config := make(map[string]*string)

		for key, value := range topicData {
			if strings.HasPrefix(key, "config.") {
				configKey := strings.TrimPrefix(key, "config.")
				config[configKey] = &value
			}
		}

		topic = &topics.Topic{
			Name:       strings.TrimPrefix(topicRedisKey, fmt.Sprintf("%s-", TopicConfigRedisKeyFmtPrefix)),
			Partitions: int32(partitions),
			Config:     config,
		}

		return nil
	}, topicRedisKey)
	if err != nil {
		return nil, nil, err
	}

	if cluster == nil {
		return nil, nil, nil
	}

	if topic == nil {
		return nil, nil, nil
	}

	return cluster, topic, nil
}
