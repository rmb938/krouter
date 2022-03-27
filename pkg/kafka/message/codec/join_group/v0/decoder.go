package v0

import (
	"time"

	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v0.Request{}

	var err error
	if msg.GroupID, err = reader.String(); err != nil {
		return nil, err
	}

	var sessionTimeoutMS int32
	if sessionTimeoutMS, err = reader.Int32(); err != nil {
		return nil, err
	}
	msg.SessionTimeout = time.Duration(sessionTimeoutMS) * time.Millisecond

	if msg.MemberID, err = reader.String(); err != nil {
		return nil, err
	}

	if msg.ProtocolType, err = reader.String(); err != nil {
		return nil, err
	}

	var protocolsLength int32
	if protocolsLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}

	for i := int32(0); i < protocolsLength; i++ {
		groupProtocol := v0.GroupProtocol{}

		if groupProtocol.Name, err = reader.String(); err != nil {
			return nil, err
		}

		if groupProtocol.Metadata, err = reader.Bytes(); err != nil {
			return nil, err
		}

		msg.Protocols = append(msg.Protocols, groupProtocol)
	}

	return msg, nil
}
