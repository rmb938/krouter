package v2

import (
	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	apiVersionv0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	apiVersionv2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups"
	implDescribeGroupsV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
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
	implListOffsetsV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v0"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, log logr.Logger, message message.Message, correlationId int32) error {
	_ = message.(*apiVersionv2.Request)

	response := &apiVersionv2.Response{}

	response.ErrCode = errors.None
	response.APIKeys = append(response.APIKeys,
		apiVersionv2.APIKey{
			Key:        produce.Key,
			MinVersion: producev7.Version,
			MaxVersion: producev7.Version,
		},
		apiVersionv2.APIKey{
			Key:        metadata.Key,
			MinVersion: metadatav8.Version,
			MaxVersion: metadatav8.Version,
		},
		apiVersionv2.APIKey{
			Key:        api_version.Key,
			MinVersion: apiVersionv0.Version,
			MaxVersion: apiVersionv2.Version,
		},
		apiVersionv2.APIKey{
			Key:        init_producer_id.Key,
			MinVersion: initProducerIDv1.Version,
			MaxVersion: initProducerIDv1.Version,
		},
		apiVersionv2.APIKey{
			Key:        find_coordinator.Key,
			MinVersion: implFindCoordinatorV2.Version,
			MaxVersion: implFindCoordinatorV2.Version,
		},
		apiVersionv2.APIKey{
			Key:        join_group.Key,
			MinVersion: implJoinGroupV4.Version,
			MaxVersion: implJoinGroupV4.Version,
		},
		apiVersionv2.APIKey{
			Key:        sync_group.Key,
			MinVersion: implSyncGroupV0.Version,
			MaxVersion: implSyncGroupV0.Version,
		},
		apiVersionv2.APIKey{
			Key:        leave_group.Key,
			MinVersion: implLeaveGroupV0.Version,
			MaxVersion: implLeaveGroupV0.Version,
		},
		apiVersionv2.APIKey{
			Key:        offset_fetch.Key,
			MinVersion: implOffsetFetchv4.Version,
			MaxVersion: implOffsetFetchv4.Version,
		},
		apiVersionv2.APIKey{
			Key:        list_offsets.Key,
			MinVersion: implListOffsetsV3.Version,
			MaxVersion: implListOffsetsV3.Version,
		},
		apiVersionv2.APIKey{
			Key:        fetch.Key,
			MinVersion: implFetchV11.Version,
			MaxVersion: implFetchV11.Version,
		},
		apiVersionv2.APIKey{
			Key:        heartbeat.Key,
			MinVersion: implHeartbeatV0.Version,
			MaxVersion: implHeartbeatV0.Version,
		},
		apiVersionv2.APIKey{
			Key:        describe_groups.Key,
			MinVersion: implDescribeGroupsV0.Version,
			MaxVersion: implDescribeGroupsV0.Version,
		},
	)

	return client.WriteMessage(response, correlationId)
}
