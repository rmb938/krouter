package v2

import (
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/handler/api_version/all"
	apiVersionv2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	_ = message.(*apiVersionv2.Request)

	response := &apiVersionv2.Response{}

	response.ErrCode = errors.None
	for apiKey, info := range all.APIKeys {
		response.APIKeys = append(response.APIKeys, apiVersionv2.APIKey{
			Key:        apiKey,
			MinVersion: info.MinVersion,
			MaxVersion: info.MaxVersion,
		})
	}

	return response, nil
}
