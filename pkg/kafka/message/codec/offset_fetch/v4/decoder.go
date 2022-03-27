package v4

import (
	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v4.Request{}

	var err error
	if msg.GroupID, err = reader.String(); err != nil {
		return nil, err
	}

	var topicsLength *int32
	if topicsLength, err = reader.NullableArrayLength(); err != nil {
		return nil, err
	}

	if topicsLength != nil {
		for i := int32(0); i < *topicsLength; i++ {
			offsetFetchTopic := v4.RequestOffsetFetchTopic{}

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
	}

	return msg, nil
}
