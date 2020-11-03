package v3

import (
	"github.com/rmb938/krouter/pkg/kafka/client"
	apiVersionv0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	apiVersionv3 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav9 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v9"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, message message.Message, correlationId int32) error {
	_ = message.(*apiVersionv3.Request)

	// api version always replies with version 0
	// https://github.com/apache/kafka/blob/2.6.0/clients/src/main/java/org/apache/kafka/common/protocol/ApiKeys.java#L160-L168
	response := &apiVersionv0.Response{}

	response.APIKeys = append(response.APIKeys,
		apiVersionv0.APIKey{
			Key:        0, // TODO: this (we are faking it for now so we can move on)
			MinVersion: 8,
			MaxVersion: 8,
		},
		apiVersionv0.APIKey{
			Key:        metadata.Key,
			MinVersion: metadatav9.Version,
			MaxVersion: metadatav9.Version,
		},
	)

	return client.WriteMessage(response, correlationId)
}
