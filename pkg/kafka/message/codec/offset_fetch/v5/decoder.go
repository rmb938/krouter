package v5

import (
	v5 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v5"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v5.Request{}

	var err error
	if msg.GroupID, err = reader.String(); err != nil {
		return nil, err
	}

	var topicsLength int32
	if topicsLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}

	for i := int32(0); i < topicsLength; i++ {
		offsetFetchTopic := v5.RequestOffsetFetchTopic{}

		if offsetFetchTopic.Name, err = reader.String(); err != nil {
			return nil, err
		}

		var partitionsLength int32
		if partitionsLength, err = reader.ArrayLength(); err != nil {
			return nil, err
		}

		for i := int32(0); i < partitionsLength; i++ {
			var partitionIndex int32
			if partitionIndex, err = reader.Int32(); err != nil {
				return nil, err
			}

			offsetFetchTopic.PartitionIndexes = append(offsetFetchTopic.PartitionIndexes, partitionIndex)
		}

		msg.Topics = append(msg.Topics, offsetFetchTopic)
	}

	return msg, nil
}
