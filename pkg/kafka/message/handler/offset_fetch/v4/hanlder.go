package v4

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(broker *logical_broker.Broker, log logr.Logger, message message.Message) (message.Message, error) {
	log = log.WithName("offset-fetch-v5-handler")

	request := message.(*v4.Request)

	log = log.WithValues("group-id", request.GroupID)

	response := &v4.Response{
		ErrCode: errors.None,
	}

	for _, requestTopic := range request.Topics {
		offsetFetchTopic := v4.ResponseOffsetFetchTopic{
			Name: requestTopic.Name,
		}

		_, topic, err := broker.GetTopic(requestTopic.Name)
		if err != nil {
			log.Error(err, "error getting topic from logical broker")
			return nil, fmt.Errorf("error getting topic from logical broker: %w", err)
		}

		for _, partitionIndex := range requestTopic.PartitionIndexes {
			offsetFetchTopicPartition := v4.ResponseOffsetFetchTopicPartition{
				PartitionIndex: partitionIndex,
				Metadata:       nil,
				ErrCode:        errors.None,
			}

			if topic == nil {
				offsetFetchTopicPartition.ErrCode = errors.UnknownTopicOrPartition
				offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, offsetFetchTopicPartition)
				continue
			}

			if partitionIndex >= topic.Partitions {
				offsetFetchTopicPartition.ErrCode = errors.UnknownTopicOrPartition
				offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, offsetFetchTopicPartition)
				continue
			}

			offset, err := broker.GetController().OffsetFetch(request.GroupID, topic.Name, partitionIndex)
			if err != nil {
				log.Error(err, "error fetching offsets to controller")
				return nil, fmt.Errorf("error fetching offsets to controller: %w", err)
			}

			if offsetFetchTopicPartition.ErrCode == errors.None {
				offsetFetchTopicPartition.CommittedOffset = offset
			}

			offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, offsetFetchTopicPartition)
		}

		response.Topics = append(response.Topics, offsetFetchTopic)
	}

	if len(request.Topics) == 0 {
		offsets, err := broker.GetController().OffsetFetchAllTopics(request.GroupID)
		if err != nil {
			log.Error(err, "error fetching offsets for all topics to controller")
			return nil, fmt.Errorf("error fetching offsets for all topics to controller: %w", err)
		}

		for topicName, partitionInfo := range offsets {
			offsetFetchTopic := v4.ResponseOffsetFetchTopic{
				Name: topicName,
			}

			for partitionIndex, offset := range partitionInfo {
				offsetFetchTopic.Partitions = append(offsetFetchTopic.Partitions, v4.ResponseOffsetFetchTopicPartition{
					PartitionIndex:  partitionIndex,
					CommittedOffset: offset,
					Metadata:        nil,
					ErrCode:         errors.None,
				})
			}

			response.Topics = append(response.Topics, offsetFetchTopic)
		}
	}

	return response, nil
}
