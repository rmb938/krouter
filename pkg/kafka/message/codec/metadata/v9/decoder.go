package v9

import (
	v9 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v9"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v9.Request{}

	topicsLength, err := reader.CompactArrayLength()
	if err != nil {
		return nil, err
	}
	if topicsLength > 0 {
		topicsLength -= 1

		for i := uint64(0); i < topicsLength; i++ {
			topicName, err := reader.CompactString()
			if err != nil {
				return nil, err
			}

			msg.Topics = append(msg.Topics, topicName)

			// tagged fields
			if _, err := reader.UVarInt(); err != nil {
				return nil, err
			}
		}
	}

	if msg.AllowAutoTopicCreation, err = reader.Bool(); err != nil {
		return nil, err
	}
	if msg.IncludeClusterAuthorizedOperations, err = reader.Bool(); err != nil {
		return nil, err
	}
	if msg.IncludeTopicAuthorizedOperations, err = reader.Bool(); err != nil {
		return nil, err
	}

	// tagged fields
	if _, err := reader.UVarInt(); err != nil {
		return nil, err
	}

	return msg, nil
}
