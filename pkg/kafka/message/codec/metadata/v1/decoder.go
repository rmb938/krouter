package v1

import (
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v1.Request{}

	msg.Topics = nil
	topicsLength, err := reader.ArrayLength()
	if err != nil {
		return nil, err
	}
	msg.Topics = make([]string, 0)
	for i := int32(0); i < topicsLength; i++ {
		topicName, err := reader.String()
		if err != nil {
			return nil, err
		}

		msg.Topics = append(msg.Topics, topicName)
	}

	return msg, nil
}
