package v0

import (
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/handler/api_version/all"
	apiVersionv0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	_ = message.(*apiVersionv0.Request)

	response := &apiVersionv0.Response{}

	response.ErrCode = errors.None

	for apiKey, info := range all.APIKeys {
		response.APIKeys = append(response.APIKeys, apiVersionv0.APIKey{
			Key:        apiKey,
			MinVersion: info.MinVersion,
			MaxVersion: info.MaxVersion,
		})
	}

	return response, nil
}
