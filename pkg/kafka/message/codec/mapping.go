package codec

import (
	"reflect"

	v0 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v0"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v2"
	fetchV11 "github.com/rmb938/krouter/pkg/kafka/message/codec/fetch/v11"
	findCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/codec/find_coordinator/v2"
	heartbeatV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/heartbeat/v0"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/codec/init_producer_id/v1"
	joinGroupV4 "github.com/rmb938/krouter/pkg/kafka/message/codec/join_group/v4"
	leaveGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/leave_group/v0"
	listOffsetsV3 "github.com/rmb938/krouter/pkg/kafka/message/codec/list_offsets/v3"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v8"
	offsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_commit/v4"
	offsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_fetch/v4"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v7"
	syncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/sync_group/v0"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	implAPIVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/fetch"
	implFetchV11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	implFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat"
	implHeartbeatV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group"
	implLeaveGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets"
	implListOffsetsV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	implMetadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit"
	implOffsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	implProducev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v0"
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
	sync_group.Key: {
		implSyncGroupV0.Version: &syncGroupV0.Decoder{},
	},
	leave_group.Key: {
		implLeaveGroupV0.Version: &leaveGroupV0.Decoder{},
	},
	offset_fetch.Key: {
		implOffsetFetchv4.Version: &offsetFetchv4.Decoder{},
	},
	list_offsets.Key: {
		implListOffsetsV3.Version: &listOffsetsV3.Decoder{},
	},
	fetch.Key: {
		implFetchV11.Version: &fetchV11.Decoder{},
	},
	heartbeat.Key: {
		implHeartbeatV0.Version: &heartbeatV0.Decoder{},
	},
	offset_commit.Key: {
		implOffsetCommitV4.Version: &offsetCommitV4.Decoder{},
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
	reflect.TypeOf(implSyncGroupV0.Response{}):       &syncGroupV0.Encoder{},
	reflect.TypeOf(implLeaveGroupV0.Response{}):      &leaveGroupV0.Encoder{},
	reflect.TypeOf(implOffsetFetchv4.Response{}):     &offsetFetchv4.Encoder{},
	reflect.TypeOf(implListOffsetsV3.Response{}):     &listOffsetsV3.Encoder{},
	reflect.TypeOf(implFetchV11.Response{}):          &fetchV11.Encoder{},
	reflect.TypeOf(implHeartbeatV0.Response{}):       &heartbeatV0.Encoder{},
	reflect.TypeOf(implOffsetCommitV4.Response{}):    &offsetCommitV4.Encoder{},
}
