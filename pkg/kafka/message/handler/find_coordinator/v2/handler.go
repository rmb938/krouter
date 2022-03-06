package v1

import (
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	_ = message.(*v2.Request)

	response := &v2.Response{}
	broker := client.Broker

	response.ThrottleDuration = 0
	response.ErrCode = errors.None
	response.ErrMessage = nil
	response.NodeID = 1
	response.Host = broker.AdvertiseListener.IP.String()
	response.Port = int32(broker.AdvertiseListener.Port)

	return client.WriteMessage(response, correlationId)
}
