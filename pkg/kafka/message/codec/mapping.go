package codec

import (
	"reflect"

	v0 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v0"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v2"
	describeConfigsV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/describe_configs/v0"
	describeGroupsV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/describe_groups/v0"
	fetchV11 "github.com/rmb938/krouter/pkg/kafka/message/codec/fetch/v11"
	fetchV4 "github.com/rmb938/krouter/pkg/kafka/message/codec/fetch/v4"
	findCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/codec/find_coordinator/v2"
	heartbeatV3 "github.com/rmb938/krouter/pkg/kafka/message/codec/heartbeat/v3"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/codec/init_producer_id/v1"
	joinGroupV5 "github.com/rmb938/krouter/pkg/kafka/message/codec/join_group/v5"
	leaveGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/codec/leave_group/v3"
	listOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/codec/list_offsets/v5"
	metadatav0 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v0"
	metadatav1 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v1"
	metadatav4 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v4"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v8"
	offsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_commit/v4"
	offsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_fetch/v4"
	producev3 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v3"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v7"
	syncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/codec/sync_group/v3"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	implAPIVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs"
	implDescribeConfigsV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups"
	implDescribeGroupsV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/fetch"
	implFetchV11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	implFetchV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	implFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat"
	implHeartbeatV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group"
	implLeaveGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets"
	implListOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	implMetadatav0 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v0"
	implMetadatav1 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v1"
	implMetadatav4 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v4"
	implMetadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit"
	implOffsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	implProducev3 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v3"
	implProducev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
	"github.com/rmb938/krouter/pkg/net/message"
)

var MessageDecoderMapping = map[int16]map[int16]message.Decoder{
	implAPIVersion.Key: {
		implAPIVersionV0.Version: &v0.Decoder{},
		implAPIVersionV2.Version: &v2.Decoder{},
	},
	metadata.Key: {
		implMetadatav8.Version: &metadatav8.Decoder{},
		implMetadatav0.Version: &metadatav0.Decoder{},
		implMetadatav1.Version: &metadatav1.Decoder{},
		implMetadatav4.Version: &metadatav4.Decoder{},
	},
	init_producer_id.Key: {
		initProducerIDV1.Version: &v1.Decoder{},
	},
	produce.Key: {
		implProducev3.Version: &producev3.Decoder{},
		implProducev7.Version: &producev7.Decoder{},
	},
	find_coordinator.Key: {
		implFindCoordinatorV2.Version: &findCoordinatorV2.Decoder{},
	},
	join_group.Key: {
		implJoinGroupV5.Version: &joinGroupV5.Decoder{},
	},
	sync_group.Key: {
		implSyncGroupV3.Version: &syncGroupV3.Decoder{},
	},
	leave_group.Key: {
		implLeaveGroupV3.Version: &leaveGroupV3.Decoder{},
	},
	offset_fetch.Key: {
		implOffsetFetchv4.Version: &offsetFetchv4.Decoder{},
	},
	list_offsets.Key: {
		implListOffsetsV5.Version: &listOffsetsV5.Decoder{},
	},
	fetch.Key: {
		implFetchV4.Version:  &fetchV4.Decoder{},
		implFetchV11.Version: &fetchV11.Decoder{},
	},
	heartbeat.Key: {
		implHeartbeatV3.Version: &heartbeatV3.Decoder{},
	},
	offset_commit.Key: {
		implOffsetCommitV4.Version: &offsetCommitV4.Decoder{},
	},
	describe_groups.Key: {
		implDescribeGroupsV0.Version: &describeGroupsV0.Decoder{},
	},
	describe_configs.Key: {
		implDescribeConfigsV0.Version: &describeConfigsV0.Decoder{},
	},
}

var MessageEncoderMapping = map[reflect.Type]message.Encoder{
	reflect.TypeOf(implAPIVersionV0.Response{}):      &v0.Encoder{},
	reflect.TypeOf(implAPIVersionV2.Response{}):      &v2.Encoder{},
	reflect.TypeOf(implMetadatav0.Response{}):        &metadatav0.Encoder{},
	reflect.TypeOf(implMetadatav1.Response{}):        &metadatav1.Encoder{},
	reflect.TypeOf(implMetadatav4.Response{}):        &metadatav4.Encoder{},
	reflect.TypeOf(implMetadatav8.Response{}):        &metadatav8.Encoder{},
	reflect.TypeOf(initProducerIDV1.Response{}):      &v1.Encoder{},
	reflect.TypeOf(implProducev3.Response{}):         &producev3.Encoder{},
	reflect.TypeOf(implProducev7.Response{}):         &producev7.Encoder{},
	reflect.TypeOf(implFindCoordinatorV2.Response{}): &findCoordinatorV2.Encoder{},
	reflect.TypeOf(implJoinGroupV5.Response{}):       &joinGroupV5.Encoder{},
	reflect.TypeOf(implSyncGroupV3.Response{}):       &syncGroupV3.Encoder{},
	reflect.TypeOf(implLeaveGroupV3.Response{}):      &leaveGroupV3.Encoder{},
	reflect.TypeOf(implOffsetFetchv4.Response{}):     &offsetFetchv4.Encoder{},
	reflect.TypeOf(implListOffsetsV5.Response{}):     &listOffsetsV5.Encoder{},
	reflect.TypeOf(implFetchV4.Response{}):           &fetchV4.Encoder{},
	reflect.TypeOf(implFetchV11.Response{}):          &fetchV11.Encoder{},
	reflect.TypeOf(implHeartbeatV3.Response{}):       &heartbeatV3.Encoder{},
	reflect.TypeOf(implOffsetCommitV4.Response{}):    &offsetCommitV4.Encoder{},
	reflect.TypeOf(implDescribeGroupsV0.Response{}):  &describeGroupsV0.Encoder{},
	reflect.TypeOf(implDescribeConfigsV0.Response{}): &describeConfigsV0.Encoder{},
}
