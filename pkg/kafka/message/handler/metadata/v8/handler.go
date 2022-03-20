package v8

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("metadata-v8-handler")

	request := message.(*metadatav8.Request)
	response := &metadatav8.Response{}

	logicalBroker := broker

	response.ThrottleDuration = 0

	uniqueClusters := make(map[*logical_broker.Cluster][]string)

	topics := request.Topics

	response.ClusterID = &logicalBroker.ClusterID
	response.ControllerID = 2147483647

	response.Brokers = append(response.Brokers, metadatav8.Brokers{
		ID:   response.ControllerID,
		Host: logicalBroker.AdvertiseListener.IP.String(),
		Port: int32(logicalBroker.AdvertiseListener.Port),
		Rack: func(s string) *string { return &s }("rack"),
	})

	if request.Topics == nil {
		for _, topic := range logicalBroker.GetTopics() {
			topics = append(topics, topic.Name)
		}
	}

	for _, topicName := range topics {
		log = log.WithValues("topic", topicName)

		cluster, _ := logicalBroker.GetTopic(topicName)

		if cluster == nil {
			log.Error(nil, "Client tried to get metadata for a topic that doesn't exist")
			response.Topics = append(response.Topics, metadatav8.Topics{
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

	uniqueBrokers := make(map[int32]struct{})
	for cluster, topics := range uniqueClusters {
		log = log.WithValues("cluster", cluster.Name)

		kafkaMetadata, err := cluster.TopicMetadata(context.TODO(), topics)
		if err != nil {
			log.Error(err, "error fetching metadata for topics from kafka")
			return nil, fmt.Errorf("error fetching metadata for topics from kafka: %w", err)
		}

		if int64(kafkaMetadata.ThrottleMillis) > response.ThrottleDuration.Milliseconds() {
			response.ThrottleDuration = time.Duration(kafkaMetadata.ThrottleMillis) * time.Millisecond
		}

		for _, broker := range kafkaMetadata.Brokers {
			if _, ok := uniqueBrokers[broker.NodeID]; !ok {
				uniqueBrokers[broker.NodeID] = struct{}{}
				response.Brokers = append(response.Brokers, metadatav8.Brokers{
					ID:   broker.NodeID,
					Host: logicalBroker.AdvertiseListener.IP.String(),
					Port: int32(logicalBroker.AdvertiseListener.Port),
					Rack: func(s string) *string { return &s }("rack"),
				})
			}
		}

		for _, topic := range kafkaMetadata.Topics {
			responseTopic := metadatav8.Topics{
				ErrCode:  errors.KafkaError(topic.ErrorCode),
				Name:     *topic.Topic,
				Internal: false,
			}

			for _, partition := range topic.Partitions {
				responsePartition := metadatav8.Partitions{
					ErrCode:         errors.KafkaError(partition.ErrorCode),
					Index:           partition.Partition,
					LeaderID:        partition.Leader,
					LeaderEpoch:     partition.LeaderEpoch,
					ReplicaNodes:    partition.Replicas,
					ISRNodes:        partition.ISR,
					OfflineReplicas: partition.OfflineReplicas,
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

		for _, broker := range kafkaMetadata.Brokers {

			if _, ok := uniqueBrokers[broker.NodeID]; !ok {
				uniqueBrokers[broker.NodeID] = struct{}{}
				response.Brokers = append(response.Brokers, metadatav8.Brokers{
					ID:   broker.NodeID,
					Host: logicalBroker.AdvertiseListener.IP.String(),
					Port: int32(logicalBroker.AdvertiseListener.Port),
					Rack: func(s string) *string { return &s }("rack"),
				})
			}
		}
	}

	return response, nil
}
