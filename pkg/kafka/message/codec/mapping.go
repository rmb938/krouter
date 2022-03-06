package codec

import (
	"reflect"

	v0 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v0"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v2"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/codec/init_producer_id/v1"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v8"
	producev5 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v5"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	implAPIVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	implMetadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	implProducev5 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v5"
	"github.com/rmb938/krouter/pkg/net/message"
)

var MessageDecoderMapping = map[int16]map[int16]message.Decoder{
	implAPIVersion.Key: {
		implAPIVersionV0.Version: &v0.Decoder{},
		implAPIVersionV2.Version: &v2.Decoder{},
	},
	metadata.Key: {
		implMetadatav8.Version: &metadatav8.Decoder{},
	},
	init_producer_id.Key: {
		initProducerIDV1.Version: &v1.Decoder{},
	},
	produce.Key: {
		implProducev5.Version: &producev5.Decoder{},
	},
}

var MessageEncoderMapping = map[reflect.Type]message.Encoder{
	reflect.TypeOf(implAPIVersionV0.Response{}): &v0.Encoder{},
	reflect.TypeOf(implAPIVersionV2.Response{}): &v2.Encoder{},
	reflect.TypeOf(implMetadatav8.Response{}):   &metadatav8.Encoder{},
	reflect.TypeOf(initProducerIDV1.Response{}): &v1.Encoder{},
	reflect.TypeOf(implProducev5.Response{}):    &producev5.Encoder{},
}
