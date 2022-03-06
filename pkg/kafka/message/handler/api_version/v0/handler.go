package v0

import (
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	apiVersionv0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	apiVersionv2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDv1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	producev5 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v5"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, message message.Message, correlationId int32) error {
	_ = message.(*apiVersionv0.Request)

	response := &apiVersionv0.Response{}

	response.ErrCode = errors.None
	response.APIKeys = append(response.APIKeys,
		apiVersionv0.APIKey{
			Key:        produce.Key,
			MinVersion: producev5.Version,
			MaxVersion: producev5.Version,
		},
		apiVersionv0.APIKey{
			Key:        metadata.Key,
			MinVersion: metadatav8.Version,
			MaxVersion: metadatav8.Version,
		},
		apiVersionv0.APIKey{
			Key:        api_version.Key,
			MinVersion: apiVersionv0.Version,
			MaxVersion: apiVersionv2.Version,
		},
		apiVersionv0.APIKey{
			Key:        init_producer_id.Key,
			MinVersion: initProducerIDv1.Version,
			MaxVersion: initProducerIDv1.Version,
		},
	)

	return client.WriteMessage(response, correlationId)
}
