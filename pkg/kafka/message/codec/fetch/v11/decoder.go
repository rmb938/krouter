package v11

import (
	"time"

	v11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v11.Request{}

	var err error
	if _, err = reader.Int32(); err != nil {
		return nil, err
	}

	var maxWaitMS int32
	if maxWaitMS, err = reader.Int32(); err != nil {
		return nil, err
	}

	msg.MaxWait = time.Duration(maxWaitMS) * time.Millisecond

	if msg.MinBytes, err = reader.Int32(); err != nil {
		return nil, err
	}

	if msg.MaxBytes, err = reader.Int32(); err != nil {
		return nil, err
	}

	if msg.IsolationLevel, err = reader.Int8(); err != nil {
		return nil, err
	}

	if msg.SessionID, err = reader.Int32(); err != nil {
		return nil, err
	}

	if msg.SessionEpoch, err = reader.Int32(); err != nil {
		return nil, err
	}

	var topicsLength int32
	if topicsLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}
	for i := int32(0); i < topicsLength; i++ {
		requestTopic := v11.FetchTopicRequest{}

		if requestTopic.Name, err = reader.String(); err != nil {
			return nil, err
		}

		var partitionsLength int32
		if partitionsLength, err = reader.ArrayLength(); err != nil {
			return nil, err
		}

		for i := int32(0); i < partitionsLength; i++ {
			requestPartition := v11.FetchPartitionRequest{}

			if requestPartition.Partition, err = reader.Int32(); err != nil {
				return nil, err
			}

			if requestPartition.CurrentLeaderEpoch, err = reader.Int32(); err != nil {
				return nil, err
			}

			if requestPartition.FetchOffset, err = reader.Int64(); err != nil {
				return nil, err
			}

			if requestPartition.LogStartOffset, err = reader.Int64(); err != nil {
				return nil, err
			}

			if requestPartition.PartitionMaxBytes, err = reader.Int32(); err != nil {
				return nil, err
			}

			requestTopic.Partitions = append(requestTopic.Partitions, requestPartition)
		}

		msg.Topics = append(msg.Topics, requestTopic)
	}

	var forgottenTopicsDataLength int32
	if forgottenTopicsDataLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}

	for i := int32(0); i < forgottenTopicsDataLength; i++ {
		forgottenTopicData := v11.FetchForgottenTopicData{}

		if forgottenTopicData.Name, err = reader.String(); err != nil {
			return nil, err
		}

		var forgottenPartitionsLength int32
		if forgottenPartitionsLength, err = reader.ArrayLength(); err != nil {
			return nil, err
		}

		for i := int32(0); i < forgottenPartitionsLength; i++ {
			var partition int32
			if partition, err = reader.Int32(); err != nil {
				return nil, err
			}

			forgottenTopicData.Partitions = append(forgottenTopicData.Partitions, partition)
		}

		msg.ForgottenTopicsData = append(msg.ForgottenTopicsData, forgottenTopicData)
	}

	if msg.RackID, err = reader.String(); err != nil {
		return nil, err
	}

	return msg, nil
}
