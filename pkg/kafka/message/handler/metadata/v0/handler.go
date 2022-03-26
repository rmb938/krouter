package v0

import (
	"context"
	"fmt"
	"math"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	metadatav0 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v0"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("metadata-v8-handler")

	request := message.(*metadatav0.Request)
	response := &metadatav0.Response{}

	logicalBroker := broker

	uniqueClusters := make(map[*logical_broker.Cluster][]string)

	topics := request.Topics

	response.Brokers = append(response.Brokers, metadatav0.Brokers{
		ID:   math.MaxInt32,
		Host: logicalBroker.AdvertiseListener.IP.String(),
		Port: int32(logicalBroker.AdvertiseListener.Port),
	})

	if request.Topics == nil {
		allTopics, err := logicalBroker.GetTopics()
		if err != nil {
			log.Error(err, "error fetching all topics brom logical broker")
			return nil, fmt.Errorf("error fetching all topics brom logical broker: %w", err)
		}
		for _, topic := range allTopics {
			topics = append(topics, topic.Name)
		}
	}

	for _, topicName := range topics {
		log = log.WithValues("topic", topicName)

		cluster, _ := logicalBroker.GetTopic(topicName)

		if cluster == nil {
			log.Error(nil, "Client tried to get metadata for a topic that doesn't exist")
			response.Topics = append(response.Topics, metadatav0.Topics{
				ErrCode: errors.UnknownTopicOrPartition,
				Name:    topicName,
			})
			continue
		}

		if _, ok := uniqueClusters[cluster]; !ok {
			uniqueClusters[cluster] = make([]string, 0)
		}

		uniqueClusters[cluster] = append(uniqueClusters[cluster], topicName)
	}

	currentBrokerId := int32(10000)
	for cluster, topics := range uniqueClusters {
		log = log.WithValues("cluster", cluster.Name)

		clusterBrokerIDMap := make(map[int32]int32)
		clusterBrokerIDMap[-1] = -1 // pre-populate -1 since it's always the same

		kafkaMetadata, err := cluster.TopicMetadata(context.TODO(), topics)
		if err != nil {
			log.Error(err, "error fetching metadata for topics from kafka")
			return nil, fmt.Errorf("error fetching metadata for topics from kafka: %w", err)
		}

		for _, broker := range kafkaMetadata.Brokers {
			response.Brokers = append(response.Brokers, metadatav0.Brokers{
				ID:   currentBrokerId,
				Host: logicalBroker.AdvertiseListener.IP.String(),
				Port: int32(logicalBroker.AdvertiseListener.Port),
			})
			clusterBrokerIDMap[broker.NodeID] = currentBrokerId
			currentBrokerId += 1
		}

		for _, topic := range kafkaMetadata.Topics {
			responseTopic := metadatav0.Topics{
				ErrCode: errors.KafkaError(topic.ErrorCode),
				Name:    *topic.Topic,
			}

			for _, partition := range topic.Partitions {
				responsePartition := metadatav0.Partitions{
					ErrCode: errors.KafkaError(partition.ErrorCode),
					Index:   partition.Partition,
				}

				responsePartition.LeaderID = clusterBrokerIDMap[partition.Leader]

				// convert the partition broker id's to ours
				for _, replicaID := range partition.Replicas {
					responsePartition.ReplicaNodes = append(responsePartition.ReplicaNodes, clusterBrokerIDMap[replicaID])
				}

				for _, isr := range partition.ISR {
					responsePartition.ISRNodes = append(responsePartition.ISRNodes, clusterBrokerIDMap[isr])
				}

				responseTopic.Partitions = append(responseTopic.Partitions, responsePartition)
			}

			response.Topics = append(response.Topics, responseTopic)
		}
	}

	// get metadata for all clusters
	for _, cluster := range logicalBroker.GetClusters() {
		if _, ok := uniqueClusters[cluster]; ok {
			// if we already got it via topics skip it
			continue
		}

		kafkaMetadata, err := cluster.TopicMetadata(context.TODO(), nil)
		if err != nil {
			log.Error(err, "error fetching metadata for brokers from kafka")
			return nil, fmt.Errorf("error fetching metadata for brokers from kafka: %w", err)
		}

		for range kafkaMetadata.Brokers {
			response.Brokers = append(response.Brokers, metadatav0.Brokers{
				ID:   currentBrokerId,
				Host: logicalBroker.AdvertiseListener.IP.String(),
				Port: int32(logicalBroker.AdvertiseListener.Port),
			})
			currentBrokerId += 1
		}
	}

	return response, nil
}
