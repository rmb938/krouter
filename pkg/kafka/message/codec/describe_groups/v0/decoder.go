package v0

import (
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v0.Request{}

	var groupLength int32
	var err error
	if groupLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}

	for i := int32(0); i < groupLength; i++ {
		var groupId string
		if groupId, err = reader.String(); err != nil {
			return nil, err
		}

		msg.Groups = append(msg.Groups, groupId)
	}

	return msg, nil
}
