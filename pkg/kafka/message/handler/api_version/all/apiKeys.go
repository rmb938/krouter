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
	implFetchV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	implFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat"
	implHeartbeatV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDv1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group"
	implLeaveGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets"
	implListOffsetsV5 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav0 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v0"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit"
	implOffsetCommitV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	producev3 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v3"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
)

type APIKey struct {
	MinVersion int16
	MaxVersion int16
}

var APIKeys = map[int16]APIKey{
	produce.Key: {
		MinVersion: producev3.Version,
		MaxVersion: producev7.Version,
	},
	metadata.Key: {
		MinVersion: metadatav0.Version,
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
		MinVersion: implJoinGroupV5.Version,
		MaxVersion: implJoinGroupV5.Version,
	},
	sync_group.Key: {
		MinVersion: implSyncGroupV3.Version,
		MaxVersion: implSyncGroupV3.Version,
	},
	leave_group.Key: {
		MinVersion: implLeaveGroupV3.Version,
		MaxVersion: implLeaveGroupV3.Version,
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
		MinVersion: implFetchV4.Version,
		MaxVersion: implFetchV11.Version,
	},
	heartbeat.Key: {
		MinVersion: implHeartbeatV3.Version,
		MaxVersion: implHeartbeatV3.Version,
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
