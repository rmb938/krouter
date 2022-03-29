package codec

import (
	"reflect"

	apiVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v0"
	apiVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/codec/api_version/v2"
	createAclsV1 "github.com/rmb938/krouter/pkg/kafka/message/codec/create_acls/v1"
	describeAclsV1 "github.com/rmb938/krouter/pkg/kafka/message/codec/describe_acls/v1"
	describeConfigsV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/describe_configs/v0"
	describeGroupsV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/describe_groups/v0"
	fetchV1 "github.com/rmb938/krouter/pkg/kafka/message/codec/fetch/v1"
	fetchV11 "github.com/rmb938/krouter/pkg/kafka/message/codec/fetch/v11"
	fetchV2 "github.com/rmb938/krouter/pkg/kafka/message/codec/fetch/v2"
	fetchV4 "github.com/rmb938/krouter/pkg/kafka/message/codec/fetch/v4"
	findCoordinatorV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/find_coordinator/v0"
	findCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/codec/find_coordinator/v2"
	heartbeatV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/heartbeat/v0"
	heartbeatV3 "github.com/rmb938/krouter/pkg/kafka/message/codec/heartbeat/v3"
	initProducerV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/init_producer_id/v0"
	initProducerV1 "github.com/rmb938/krouter/pkg/kafka/message/codec/init_producer_id/v1"
	joinGroupV5 "github.com/rmb938/krouter/pkg/kafka/message/codec/join_group/v5"
	leaveGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/leave_group/v0"
	leaveGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/codec/leave_group/v3"
	listOffsetsV1 "github.com/rmb938/krouter/pkg/kafka/message/codec/list_offsets/v1"
	listOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/codec/list_offsets/v5"
	metadatav0 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v0"
	metadatav1 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v1"
	metadatav4 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v4"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/codec/metadata/v8"
	offsetCommitV1 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_commit/v1"
	offsetCommitV2 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_commit/v2"
	offsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_commit/v4"
	offsetFetchv1 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_fetch/v1"
	offsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/codec/offset_fetch/v4"
	producev1 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v1"
	producev2 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v2"
	producev3 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v3"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/codec/produce/v7"
	syncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/codec/sync_group/v0"
	syncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/codec/sync_group/v3"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	implAPIVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/create_acls"
	implCreateAclsV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/create_acls/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_acls"
	implDescribeAclsV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_acls/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs"
	implDescribeConfigsV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups"
	implDescribeGroupsV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/fetch"
	implFetchV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v1"
	implFetchV11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	implFetchV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v2"
	implFetchV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	implFindCoordinatorV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v0"
	implFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat"
	implHeartbeatV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v0"
	implHeartbeatV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v0"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group"
	implLeaveGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v0"
	implLeaveGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets"
	implListOffsetsV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v1"
	implListOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	implMetadatav0 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v0"
	implMetadatav1 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v1"
	implMetadatav4 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v4"
	implMetadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit"
	implOffsetCommitV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v1"
	implOffsetCommitV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v2"
	implOffsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv1 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v1"
	implOffsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	implProducev1 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v1"
	implProducev2 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v2"
	implProducev3 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v3"
	implProducev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v0"
	implSyncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
	"github.com/rmb938/krouter/pkg/net/message"
)

var MessageDecoderMapping = map[int16]map[int16]message.Decoder{
	implAPIVersion.Key: {
		implAPIVersionV0.Version: &apiVersionV0.Decoder{},
		implAPIVersionV2.Version: &apiVersionV2.Decoder{},
	},
	metadata.Key: {
		implMetadatav8.Version: &metadatav8.Decoder{},
		implMetadatav0.Version: &metadatav0.Decoder{},
		implMetadatav1.Version: &metadatav1.Decoder{},
		implMetadatav4.Version: &metadatav4.Decoder{},
	},
	init_producer_id.Key: {
		initProducerIDV0.Version: &initProducerV0.Decoder{},
		initProducerIDV1.Version: &initProducerV1.Decoder{},
	},
	produce.Key: {
		implProducev1.Version: &producev1.Decoder{},
		implProducev2.Version: &producev2.Decoder{},
		implProducev3.Version: &producev3.Decoder{},
		implProducev7.Version: &producev7.Decoder{},
	},
	find_coordinator.Key: {
		implFindCoordinatorV0.Version: &findCoordinatorV0.Decoder{},
		implFindCoordinatorV2.Version: &findCoordinatorV2.Decoder{},
	},
	join_group.Key: {
		implJoinGroupV5.Version: &joinGroupV5.Decoder{},
	},
	sync_group.Key: {
		implSyncGroupV0.Version: &syncGroupV0.Decoder{},
		implSyncGroupV3.Version: &syncGroupV3.Decoder{},
	},
	leave_group.Key: {
		implLeaveGroupV0.Version: &leaveGroupV0.Decoder{},
		implLeaveGroupV3.Version: &leaveGroupV3.Decoder{},
	},
	offset_fetch.Key: {
		implOffsetFetchv1.Version: &offsetFetchv1.Decoder{},
		implOffsetFetchv4.Version: &offsetFetchv4.Decoder{},
	},
	list_offsets.Key: {
		implListOffsetsV1.Version: &listOffsetsV1.Decoder{},
		implListOffsetsV5.Version: &listOffsetsV5.Decoder{},
	},
	fetch.Key: {
		implFetchV1.Version:  &fetchV1.Decoder{},
		implFetchV2.Version:  &fetchV2.Decoder{},
		implFetchV4.Version:  &fetchV4.Decoder{},
		implFetchV11.Version: &fetchV11.Decoder{},
	},
	heartbeat.Key: {
		implHeartbeatV0.Version: &heartbeatV0.Decoder{},
		implHeartbeatV3.Version: &heartbeatV3.Decoder{},
	},
	offset_commit.Key: {
		implOffsetCommitV1.Version: &offsetCommitV1.Decoder{},
		implOffsetCommitV2.Version: &offsetCommitV2.Decoder{},
		implOffsetCommitV4.Version: &offsetCommitV4.Decoder{},
	},
	describe_groups.Key: {
		implDescribeGroupsV0.Version: &describeGroupsV0.Decoder{},
	},
	describe_configs.Key: {
		implDescribeConfigsV0.Version: &describeConfigsV0.Decoder{},
	},
	describe_acls.Key: {
		implDescribeAclsV1.Version: &describeAclsV1.Decoder{},
	},
	create_acls.Key: {
		implCreateAclsV1.Version: &createAclsV1.Decoder{},
	},
}

var MessageEncoderMapping = map[reflect.Type]message.Encoder{
	reflect.TypeOf(implAPIVersionV0.Response{}):      &apiVersionV0.Encoder{},
	reflect.TypeOf(implAPIVersionV2.Response{}):      &apiVersionV2.Encoder{},
	reflect.TypeOf(implMetadatav0.Response{}):        &metadatav0.Encoder{},
	reflect.TypeOf(implMetadatav1.Response{}):        &metadatav1.Encoder{},
	reflect.TypeOf(implMetadatav4.Response{}):        &metadatav4.Encoder{},
	reflect.TypeOf(implMetadatav8.Response{}):        &metadatav8.Encoder{},
	reflect.TypeOf(initProducerIDV0.Response{}):      &initProducerV0.Encoder{},
	reflect.TypeOf(initProducerIDV1.Response{}):      &initProducerV1.Encoder{},
	reflect.TypeOf(implProducev1.Response{}):         &producev1.Encoder{},
	reflect.TypeOf(implProducev2.Response{}):         &producev2.Encoder{},
	reflect.TypeOf(implProducev3.Response{}):         &producev3.Encoder{},
	reflect.TypeOf(implProducev7.Response{}):         &producev7.Encoder{},
	reflect.TypeOf(implFindCoordinatorV0.Response{}): &findCoordinatorV0.Encoder{},
	reflect.TypeOf(implFindCoordinatorV2.Response{}): &findCoordinatorV2.Encoder{},
	reflect.TypeOf(implJoinGroupV5.Response{}):       &joinGroupV5.Encoder{},
	reflect.TypeOf(implSyncGroupV0.Response{}):       &syncGroupV0.Encoder{},
	reflect.TypeOf(implSyncGroupV3.Response{}):       &syncGroupV3.Encoder{},
	reflect.TypeOf(implLeaveGroupV3.Response{}):      &leaveGroupV3.Encoder{},
	reflect.TypeOf(implOffsetFetchv1.Response{}):     &offsetFetchv1.Encoder{},
	reflect.TypeOf(implOffsetFetchv4.Response{}):     &offsetFetchv4.Encoder{},
	reflect.TypeOf(implListOffsetsV1.Response{}):     &listOffsetsV1.Encoder{},
	reflect.TypeOf(implListOffsetsV5.Response{}):     &listOffsetsV5.Encoder{},
	reflect.TypeOf(implFetchV1.Response{}):           &fetchV1.Encoder{},
	reflect.TypeOf(implFetchV2.Response{}):           &fetchV2.Encoder{},
	reflect.TypeOf(implFetchV4.Response{}):           &fetchV4.Encoder{},
	reflect.TypeOf(implFetchV11.Response{}):          &fetchV11.Encoder{},
	reflect.TypeOf(implHeartbeatV0.Response{}):       &heartbeatV0.Encoder{},
	reflect.TypeOf(implHeartbeatV3.Response{}):       &heartbeatV3.Encoder{},
	reflect.TypeOf(implOffsetCommitV1.Response{}):    &offsetCommitV1.Encoder{},
	reflect.TypeOf(implOffsetCommitV2.Response{}):    &offsetCommitV2.Encoder{},
	reflect.TypeOf(implOffsetCommitV4.Response{}):    &offsetCommitV4.Encoder{},
	reflect.TypeOf(implDescribeGroupsV0.Response{}):  &describeGroupsV0.Encoder{},
	reflect.TypeOf(implDescribeConfigsV0.Response{}): &describeConfigsV0.Encoder{},
	reflect.TypeOf(implDescribeAclsV1.Response{}):    &describeAclsV1.Encoder{},
	reflect.TypeOf(implCreateAclsV1.Response{}):      &createAclsV1.Encoder{},
}
