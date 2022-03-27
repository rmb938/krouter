package v1

import (
	"time"

	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {
	msg := &v1.Request{}

	var err error
	if msg.GroupID, err = reader.String(); err != nil {
		return nil, err
	}

	if msg.GenerationID, err = reader.Int32(); err != nil {
		return nil, err
	}

	if msg.MemberID, err = reader.String(); err != nil {
		return nil, err
	}

	var topicsLength int32
	if topicsLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}

	for i := int32(0); i < topicsLength; i++ {
		topicRequest := v1.OffsetCommitTopicRequest{}

		if topicRequest.Name, err = reader.String(); err != nil {
			return nil, err
		}

		var partitionsLength int32
		if partitionsLength, err = reader.ArrayLength(); err != nil {
			return nil, err
		}

		for i := int32(0); i < partitionsLength; i++ {
			partitionRequest := v1.OffsetCommitPartitionRequest{}

			if partitionRequest.PartitionIndex, err = reader.Int32(); err != nil {
				return nil, err
			}

			if partitionRequest.CommittedOffset, err = reader.Int64(); err != nil {
				return nil, err
			}

			var committedTimestamp int64
			if committedTimestamp, err = reader.Int64(); err != nil {
				return nil, err
			}
			partitionRequest.CommittedTimestamp = time.UnixMilli(committedTimestamp)

			if partitionRequest.CommittedMetadata, err = reader.NullableString(); err != nil {
				return nil, err
			}

			topicRequest.Partitions = append(topicRequest.Partitions, partitionRequest)
		}

		msg.Topics = append(msg.Topics, topicRequest)
	}

	return msg, nil
}
