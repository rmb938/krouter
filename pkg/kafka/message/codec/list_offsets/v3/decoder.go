package v3

import (
	"time"

	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/list_offsets/v3"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v3.Request{}

	var err error
	if msg.ReplicaID, err = reader.Int32(); err != nil {
		return nil, err
	}

	if msg.IsolationLevel, err = reader.Int8(); err != nil {
		return nil, err
	}

	var topicsLength int32
	if topicsLength, err = reader.Int32(); err != nil {
		return nil, err
	}
	for i := int32(0); i < topicsLength; i++ {
		topicRequest := v3.ListOffsetsTopicRequest{}

		if topicRequest.Name, err = reader.String(); err != nil {
			return nil, err
		}

		var partitionsLength int32
		if partitionsLength, err = reader.Int32(); err != nil {
			return nil, err
		}
		for i := int32(0); i < partitionsLength; i++ {
			partitionRequest := v3.ListOffsetsPartitionRequest{}

			if partitionRequest.PartitionIndex, err = reader.Int32(); err != nil {
				return nil, err
			}

			var partitionTimestamp int64
			if partitionTimestamp, err = reader.Int64(); err != nil {
				return nil, err
			}

			partitionRequest.Timestamp = time.UnixMilli(partitionTimestamp)

			topicRequest.Partitions = append(topicRequest.Partitions, partitionRequest)
		}

		msg.Topics = append(msg.Topics, topicRequest)
	}

	return msg, nil
}
