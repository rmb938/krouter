package handler

import (
	v3 "github.com/rmb938/krouter/pkg/kafka/message/handler/api_version/v3"
	v9 "github.com/rmb938/krouter/pkg/kafka/message/handler/metadata/v9"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav9 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v9"
	"github.com/rmb938/krouter/pkg/net/message/handler"
)

var MessageHandlerMapping = map[int16]map[int16]handler.MessageHandler{
	implAPIVersion.Key: {
		implAPIVersionV3.Version: &v3.Handler{},
	},
	metadata.Key: {
		metadatav9.Version: &v9.Handler{},
	},
}
