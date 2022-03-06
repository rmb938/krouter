package codec

import (
	"reflect"

	v0 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v0"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v2"
	findCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/codec/find_coordinator/v2"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/codec/init_producer_id/v1"
	joinGroupV4 "github.com/rmb938/krouter/pkg/kafka/message/codec/join_group/v4"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v8"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v7"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	implAPIVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	implFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	implMetadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	implProducev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
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
		implProducev7.Version: &producev7.Decoder{},
	},
	find_coordinator.Key: {
		implFindCoordinatorV2.Version: &findCoordinatorV2.Decoder{},
	},
	join_group.Key: {
		implJoinGroupV4.Version: &joinGroupV4.Decoder{},
	},
}

var MessageEncoderMapping = map[reflect.Type]message.Encoder{
	reflect.TypeOf(implAPIVersionV0.Response{}):      &v0.Encoder{},
	reflect.TypeOf(implAPIVersionV2.Response{}):      &v2.Encoder{},
	reflect.TypeOf(implMetadatav8.Response{}):        &metadatav8.Encoder{},
	reflect.TypeOf(initProducerIDV1.Response{}):      &v1.Encoder{},
	reflect.TypeOf(implProducev7.Response{}):         &producev7.Encoder{},
	reflect.TypeOf(implFindCoordinatorV2.Response{}): &findCoordinatorV2.Encoder{},
	reflect.TypeOf(implJoinGroupV4.Response{}):       &joinGroupV4.Encoder{},
}
