package v4

import (
	"fmt"

	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v4"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v4.Request{}

	msg.Topics = nil
	topicsLength, err := reader.NullableArrayLength()
	if err != nil {
		return nil, err
	}
	if topicsLength != nil {
		msg.Topics = make([]string, 0)
		for i := int32(0); i < *topicsLength; i++ {
			topicName, err := reader.String()
			if err != nil {
				return nil, err
			}

			msg.Topics = append(msg.Topics, topicName)
		}
	}

	if msg.AllowAutoTopicCreation, err = reader.Bool(); err != nil {
		return nil, fmt.Errorf("error bool for auto topic creation: %w", err)
	}

	return msg, nil
}
