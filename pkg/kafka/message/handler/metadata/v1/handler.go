package v1

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	metadatav1 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v1"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.LogicalBroker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("metadata-v1-handler")

	request := message.(*metadatav1.Request)
	response := &metadatav1.Response{}

	logicalBroker := broker

	uniqueClusters := make(map[*logical_broker.Cluster][]string)

	topics := request.Topics

	for _, rb := range logicalBroker.GetRegisteredBrokers() {
		response.Brokers = append(response.Brokers, metadatav1.Brokers{
			ID:   rb.ID,
			Host: rb.Endpoint.Host,
			Port: rb.Endpoint.Port,
			Rack: &rb.Rack,
		})

		// controller is always the biggest
		//  this doesn't really matter tbh
		if rb.ID > response.ControllerID {
			response.ControllerID = rb.ID
		}
	}

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

		cluster := logicalBroker.GetClusterByTopic(topicName)

		if cluster == nil {
			log.Error(nil, "Client tried to get metadata for a topic that doesn't exist")
			response.Topics = append(response.Topics, metadatav1.Topics{
				ErrCode:  errors.UnknownTopicOrPartition,
				Name:     topicName,
				Internal: false,
			})
			continue
		}

		if _, ok := uniqueClusters[cluster]; !ok {
			uniqueClusters[cluster] = make([]string, 0)
		}

		uniqueClusters[cluster] = append(uniqueClusters[cluster], topicName)
	}

	for cluster, topics := range uniqueClusters {
		log = log.WithValues("cluster", cluster.Name)

		kafkaMetadata, err := cluster.TopicMetadata(context.TODO(), topics)
		if err != nil {
			log.Error(err, "error fetching metadata for topics from kafka")
			return nil, fmt.Errorf("error fetching metadata for topics from kafka: %w", err)
		}

		for _, topic := range kafkaMetadata.Topics {
			responseTopic := metadatav1.Topics{
				ErrCode:  errors.KafkaError(topic.ErrorCode),
				Name:     *topic.Topic,
				Internal: false,
			}

			for _, partition := range topic.Partitions {
				responsePartition := metadatav1.Partitions{
					ErrCode:      errors.KafkaError(partition.ErrorCode),
					Index:        partition.Partition,
					LeaderID:     partition.Leader,
					ReplicaNodes: partition.Replicas,
					ISRNodes:     partition.ISR,
				}

				responseTopic.Partitions = append(responseTopic.Partitions, responsePartition)
			}

			response.Topics = append(response.Topics, responseTopic)
		}
	}

	return response, nil
}
