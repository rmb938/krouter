package handler

import (
	handlerAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/api_version/v0"
	handlerAPIVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/handler/api_version/v2"
	handlerDescribeConfigsV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/describe_configs/v0"
	handlerDescribeGroupsV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/describe_groups/v0"
	handlerFetchV1 "github.com/rmb938/krouter/pkg/kafka/message/handler/fetch/v1"
	handlerFetchV11 "github.com/rmb938/krouter/pkg/kafka/message/handler/fetch/v11"
	handlerFetchV2 "github.com/rmb938/krouter/pkg/kafka/message/handler/fetch/v2"
	handlerFetchV4 "github.com/rmb938/krouter/pkg/kafka/message/handler/fetch/v4"
	handlerFindCoordinatorV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/find_coordinator/v0"
	handlerFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/handler/find_coordinator/v2"
	handlerHearbeatV3 "github.com/rmb938/krouter/pkg/kafka/message/handler/heartbeat/v3"
	handlerInitProducerIDV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/init_producer_id/v0"
	handlerInitProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/handler/init_producer_id/v1"
	handlerJoinGroupV5 "github.com/rmb938/krouter/pkg/kafka/message/handler/join_group/v5"
	handlerLeaveGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/handler/leave_group/v3"
	handlerListOffsetsV1 "github.com/rmb938/krouter/pkg/kafka/message/handler/list_offsets/v1"
	handlerListOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/handler/list_offsets/v5"
	handlerMetadataV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/metadata/v0"
	handlerMetadataV1 "github.com/rmb938/krouter/pkg/kafka/message/handler/metadata/v1"
	handlerMetadataV4 "github.com/rmb938/krouter/pkg/kafka/message/handler/metadata/v4"
	v8 "github.com/rmb938/krouter/pkg/kafka/message/handler/metadata/v8"
	handlerOffsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/handler/offset_commit/v4"
	handlerOffsetFetchV4 "github.com/rmb938/krouter/pkg/kafka/message/handler/offset_fetch/v4"
	handlerProduceV1 "github.com/rmb938/krouter/pkg/kafka/message/handler/produce/v1"
	handlerProduceV2 "github.com/rmb938/krouter/pkg/kafka/message/handler/produce/v2"
	handlerProduceV3 "github.com/rmb938/krouter/pkg/kafka/message/handler/produce/v3"
	handlerProduceV7 "github.com/rmb938/krouter/pkg/kafka/message/handler/produce/v7"
	handlerSyncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/handler/sync_group/v3"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	implAPIVersionV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
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
	implHeartbeatV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v0"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group"
	implLeaveGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets"
	implListOffsetsV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v1"
	implListOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	implMetadataV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v0"
	implMetadataV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v1"
	implMetadataV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v4"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit"
	implOffsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	producev1 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v1"
	producev2 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v2"
	producev3 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v3"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
	"github.com/rmb938/krouter/pkg/net/message/handler"
)

var MessageHandlerMapping = map[int16]map[int16]handler.MessageHandler{
	implAPIVersion.Key: {
		implAPIVersionV0.Version: &handlerAPIVersionV0.Handler{},
		implAPIVersionV2.Version: &handlerAPIVersionV2.Handler{},
	},
	metadata.Key: {
		implMetadataV0.Version: &handlerMetadataV0.Handler{},
		implMetadataV1.Version: &handlerMetadataV1.Handler{},
		implMetadataV4.Version: &handlerMetadataV4.Handler{},
		metadatav8.Version:     &v8.Handler{},
	},
	init_producer_id.Key: {
		initProducerIDV0.Version: &handlerInitProducerIDV0.Handler{},
		initProducerIDV1.Version: &handlerInitProducerIDV1.Handler{},
	},
	produce.Key: {
		producev1.Version: &handlerProduceV1.Handler{},
		producev2.Version: &handlerProduceV2.Handler{},
		producev3.Version: &handlerProduceV3.Handler{},
		producev7.Version: &handlerProduceV7.Handler{},
	},
	find_coordinator.Key: {
		implFindCoordinatorV0.Version: &handlerFindCoordinatorV0.Handler{},
		implFindCoordinatorV2.Version: &handlerFindCoordinatorV2.Handler{},
	},
	join_group.Key: {
		implJoinGroupV5.Version: &handlerJoinGroupV5.Handler{},
	},
	sync_group.Key: {
		implSyncGroupV3.Version: &handlerSyncGroupV3.Handler{},
	},
	leave_group.Key: {
		implLeaveGroupV3.Version: &handlerLeaveGroupV3.Handler{},
	},
	offset_fetch.Key: {
		implOffsetFetchv4.Version: &handlerOffsetFetchV4.Handler{},
	},
	list_offsets.Key: {
		implListOffsetsV1.Version: &handlerListOffsetsV1.Handler{},
		implListOffsetsV5.Version: &handlerListOffsetsV5.Handler{},
	},
	fetch.Key: {
		implFetchV1.Version:  &handlerFetchV1.Handler{},
		implFetchV2.Version:  &handlerFetchV2.Handler{},
		implFetchV4.Version:  &handlerFetchV4.Handler{},
		implFetchV11.Version: &handlerFetchV11.Handler{},
	},
	heartbeat.Key: {
		implHeartbeatV3.Version: &handlerHearbeatV3.Handler{},
	},
	offset_commit.Key: {
		implOffsetCommitV4.Version: &handlerOffsetCommitV4.Handler{},
	},
	describe_groups.Key: {
		implDescribeGroupsV0.Version: &handlerDescribeGroupsV0.Handler{},
	},
	describe_configs.Key: {
		implDescribeConfigsV0.Version: &handlerDescribeConfigsV0.Handler{},
	},
}
