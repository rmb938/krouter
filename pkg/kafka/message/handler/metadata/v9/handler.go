package v9

import (
	"context"
	"fmt"

	"github.com/rmb938/krouter/pkg/kafka/client"
	metadatav9 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v9"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, message message.Message, correlationId int32) error {
	_ = message.(*metadatav9.Request)

	response := &metadatav9.Response{}

	brokers, err := client.Broker.GetBrokers(context.TODO())
	if err != nil {
		return fmt.Errorf("error getting brokers: %w", err)
	}

	for _, broker := range brokers {
		response.Brokers = append(response.Brokers, metadatav9.Brokers{
			ID:   broker.ID,
			Host: broker.Host,
			Port: broker.Port,
			Rack: broker.Rack,
		})
	}

	controllerID, err := client.Broker.GetControllerID(context.TODO())
	if err != nil {
		return fmt.Errorf("error getting controller id: %w", err)
	}

	if controllerID == nil {
		controllerID = func(i int32) *int32 { return &i }(-1)
	}

	response.ClusterID = client.Broker.ClusterID
	response.ControllerID = *controllerID
	response.Topics = nil
	response.ClusterAuthorizedOperations = 0

	return client.WriteMessage(response, correlationId)
}
