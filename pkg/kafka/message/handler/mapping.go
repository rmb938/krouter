package handler

import (
	v0 "github.com/rmb938/krouter/pkg/kafka/message/handler/api_version/v0"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/handler/init_producer_id/v1"
	v8 "github.com/rmb938/krouter/pkg/kafka/message/handler/metadata/v8"
	handlerProduceV5 "github.com/rmb938/krouter/pkg/kafka/message/handler/produce/v5"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	producev5 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v5"
	"github.com/rmb938/krouter/pkg/net/message/handler"
)

var MessageHandlerMapping = map[int16]map[int16]handler.MessageHandler{
	implAPIVersion.Key: {
		implAPIVersionV0.Version: &v0.Handler{},
	},
	metadata.Key: {
		metadatav8.Version: &v8.Handler{},
	},
	init_producer_id.Key: {
		initProducerIDV1.Version: &v1.Handler{},
	},
	produce.Key: {
		producev5.Version: &handlerProduceV5.Handler{},
	},
}
