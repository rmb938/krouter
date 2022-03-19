package all

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	apiVersionv0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	apiVersionv2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs"
	implDescribeConfigsV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups"
	implDescribeGroupsV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/fetch"
	implFetchV11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	implFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat"
	implHeartbeatV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDv1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group"
	implLeaveGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets"
	implListOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit"
	implOffsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v0"
)

type APIKey struct {
	MinVersion int16
	MaxVersion int16
}

var APIKeys = map[int16]APIKey{
	produce.Key: {
		MinVersion: producev7.Version,
		MaxVersion: producev7.Version,
	},
	metadata.Key: {
		MinVersion: metadatav8.Version,
		MaxVersion: metadatav8.Version,
	},
	api_version.Key: {
		MinVersion: apiVersionv0.Version,
		MaxVersion: apiVersionv2.Version,
	},
	init_producer_id.Key: {
		MinVersion: initProducerIDv1.Version,
		MaxVersion: initProducerIDv1.Version,
	},
	find_coordinator.Key: {
		MinVersion: implFindCoordinatorV2.Version,
		MaxVersion: implFindCoordinatorV2.Version,
	},
	join_group.Key: {
		MinVersion: implJoinGroupV4.Version,
		MaxVersion: implJoinGroupV4.Version,
	},
	sync_group.Key: {
		MinVersion: implSyncGroupV0.Version,
		MaxVersion: implSyncGroupV0.Version,
	},
	leave_group.Key: {
		MinVersion: implLeaveGroupV0.Version,
		MaxVersion: implLeaveGroupV0.Version,
	},
	offset_fetch.Key: {
		MinVersion: implOffsetFetchv4.Version,
		MaxVersion: implOffsetFetchv4.Version,
	},
	list_offsets.Key: {
		MinVersion: implListOffsetsV5.Version,
		MaxVersion: implListOffsetsV5.Version,
	},
	fetch.Key: {
		MinVersion: implFetchV11.Version,
		MaxVersion: implFetchV11.Version,
	},
	heartbeat.Key: {
		MinVersion: implHeartbeatV0.Version,
		MaxVersion: implHeartbeatV0.Version,
	},
	offset_commit.Key: {
		MinVersion: implOffsetCommitV4.Version,
		MaxVersion: implOffsetCommitV4.Version,
	},
	describe_groups.Key: {
		MinVersion: implDescribeGroupsV0.Version,
		MaxVersion: implDescribeGroupsV0.Version,
	},
	describe_configs.Key: {
		MinVersion: implDescribeConfigsV0.Version,
		MaxVersion: implDescribeConfigsV0.Version,
	},
}
