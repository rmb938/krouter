package codec

import (
	"reflect"

	v0 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v0"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v3"
	v9 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v9"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	implAPIVersionV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav9 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v9"
	"github.com/rmb938/krouter/pkg/net/message"
)

var MessageDecoderMapping = map[int16]map[int16]message.Decoder{
	implAPIVersion.Key: {
		implAPIVersionV3.Version: &v3.Decoder{},
	},
	metadata.Key: {
		metadatav9.Version: &v9.Decoder{},
	},
}

var MessageEncoderMapping = map[reflect.Type]message.Encoder{
	reflect.TypeOf(implAPIVersionV0.Response{}): &v0.Encoder{},
	reflect.TypeOf(metadatav9.Response{}):       &v9.Encoder{},
}
