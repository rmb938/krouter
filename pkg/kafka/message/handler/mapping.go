package handler

import (
	v0 "github.com/rmb938/krouter/pkg/kafka/message/handler/api_version/v0"
	handlerFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/handler/find_coordinator/v2"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/handler/init_producer_id/v1"
	handlerJoinGroupV4 "github.com/rmb938/krouter/pkg/kafka/message/handler/join_group/v4"
	handlerLeaveGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/leave_group/v0"
	handlerListOffsetsV3 "github.com/rmb938/krouter/pkg/kafka/message/handler/list_offsets/v3"
	v8 "github.com/rmb938/krouter/pkg/kafka/message/handler/metadata/v8"
	handlerOffsetFetchV5 "github.com/rmb938/krouter/pkg/kafka/message/handler/offset_fetch/v5"
	handlerProduceV7 "github.com/rmb938/krouter/pkg/kafka/message/handler/produce/v7"
	handlerSyncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/handler/sync_group/v0"
	implAPIVersion "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	implAPIVersionV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	implFindCoordinatorV2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	initProducerIDV1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/join_group"
	implJoinGroupV4 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v4"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group"
	implLeaveGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets"
	implListOffsetsV3 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v3"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	metadatav8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch"
	implOffsetFetchv5 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v5"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	producev7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group"
	implSyncGroupV0 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v0"
	"github.com/rmb938/krouter/pkg/net/message/handler"
)

var MessageHandlerMapping = map[int16]map[int16]handler.MessageHandler{
	implAPIVersion.Key: {
		implAPIVersionV0.Version: &v0.Handler{},
	},
	metadata.Key: {
		metadatav8.Version: &v8.Handler{},
	},
	init_producer_id.Key: {
		initProducerIDV1.Version: &v1.Handler{},
	},
	produce.Key: {
		producev7.Version: &handlerProduceV7.Handler{},
	},
	find_coordinator.Key: {
		implFindCoordinatorV2.Version: &handlerFindCoordinatorV2.Handler{},
	},
	join_group.Key: {
		implJoinGroupV4.Version: &handlerJoinGroupV4.Handler{},
	},
	sync_group.Key: {
		implSyncGroupV0.Version: &handlerSyncGroupV0.Handler{},
	},
	leave_group.Key: {
		implLeaveGroupV0.Version: &handlerLeaveGroupV0.Handler{},
	},
	offset_fetch.Key: {
		implOffsetFetchv5.Version: &handlerOffsetFetchV5.Handler{},
	},
	list_offsets.Key: {
		implListOffsetsV3.Version: &handlerListOffsetsV3.Handler{},
	},
}
