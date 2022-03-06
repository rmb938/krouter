package v8

import (
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, message message.Message, correlationId int32) error {
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
		// TODO: don't just blindly return this, pull from known topics

		response.Topics = append(response.Topics,
			metadatav8.Topics{
				ErrCode:  errors.None,
				Name:     topicName,
				Internal: false,
				Partitions: []metadatav8.Partitions{
					{
						Index:           0,
						LeaderID:        1,
						ReplicaNodes:    []int32{1},
						ISRNodes:        []int32{1},
						OfflineReplicas: []int32{},
					},
				},
			},
		)
	}

	return client.WriteMessage(response, correlationId)
}
