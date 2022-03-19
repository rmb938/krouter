package v3

import (
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v3"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v3.Request{}

	var err error
	if msg.GroupID, err = reader.String(); err != nil {
		return nil, err
	}

	var membersLength int32
	if membersLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}
	for i := int32(0); i < membersLength; i++ {
		requestMember := v3.RequestMember{}

		if requestMember.MemberID, err = reader.String(); err != nil {
			return nil, err
		}

		if requestMember.GroupInstanceID, err = reader.NullableString(); err != nil {
			return nil, err
		}

		msg.Members = append(msg.Members, requestMember)
	}

	return msg, nil
}
