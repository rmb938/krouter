package v8

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	log = log.WithName("metadata-v8-handler")

	request := message.(*metadatav8.Request)

	response := &metadatav8.Response{}

	response.ThrottleDuration = 0
	broker := client.Broker

	response.Brokers = append(response.Brokers, metadatav8.Brokers{
		ID:   1,
		Host: broker.AdvertiseListener.IP.String(),
		Port: int32(broker.AdvertiseListener.Port),
		Rack: func(s string) *string { return &s }("rack"),
	})

	response.ClusterID = &client.Broker.ClusterID
	response.ControllerID = 1

	topicNames := request.Topics

	// if requested topics is empty return all topics
	if len(topicNames) == 0 {
		for _, topic := range client.Broker.GetTopics() {
			topicNames = append(topicNames, topic.Name)
		}
	}

	for _, topicName := range topicNames {
		log = log.WithValues("topic", topicName)

		responseTopic := metadatav8.Topics{
			ErrCode:  errors.UnknownTopicOrPartition,
			Name:     topicName,
			Internal: false,
		}

		cluster, topic := client.Broker.GetTopic(topicName)
		if topic != nil {
			responseTopic.ErrCode = errors.None

			for i := int32(0); i < topic.Partitions; i++ {
				log = log.WithValues("partition", i)

				partition := metadatav8.Partitions{
					Index:           i,
					LeaderEpoch:     0, // TODO: this?
					ReplicaNodes:    []int32{},
					ISRNodes:        []int32{},
					OfflineReplicas: []int32{},
				}

				leaderID, err := cluster.TopicLeader(context.TODO(), topic.Name, i)
				if err != nil {
					log.Error(err, "Error finding topic partition leader")
					partition.ErrCode = errors.UnknownServerError
					responseTopic.Partitions = append(responseTopic.Partitions, partition)
					continue
				}

				if leaderID != -1 {
					partition.LeaderID = 1
					partition.ReplicaNodes = append(partition.ReplicaNodes, 1)
					partition.ISRNodes = append(partition.ISRNodes, 1)
				} else {
					// TODO: schedule metadata refresh
					partition.ErrCode = errors.LeaderNotAvailable
				}

				responseTopic.Partitions = append(responseTopic.Partitions, partition)
			}
		} else {
			log.Error(nil, "Client tried to get metadata for a topic that doesn't exist")
		}

		response.Topics = append(response.Topics, responseTopic)
	}

	return client.WriteMessage(response, correlationId)
}
