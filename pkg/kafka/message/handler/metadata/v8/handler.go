package v8

import (
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

	for _, topicName := range request.Topics {
		log = log.WithValues("topic", topicName)

		responseTopic := metadatav8.Topics{
			ErrCode:  errors.UnknownTopicOrPartition,
			Name:     topicName,
			Internal: false,
		}

		_, topic := client.Broker.GetTopic(topicName)
		if topic != nil {
			responseTopic.ErrCode = errors.None

			for i := int32(0); i < topic.Partitions; i++ {
				responseTopic.Partitions = append(responseTopic.Partitions,
					metadatav8.Partitions{
						Index:           i,
						LeaderID:        1,
						LeaderEpoch:     0, // TODO: this?
						ReplicaNodes:    []int32{1},
						ISRNodes:        []int32{1},
						OfflineReplicas: []int32{},
					})
			}
		} else {
			log.V(1).Info("Client tried to get metadata for a topic that doesn't exist")
		}

		response.Topics = append(response.Topics, responseTopic)
	}

	return client.WriteMessage(response, correlationId)
}
