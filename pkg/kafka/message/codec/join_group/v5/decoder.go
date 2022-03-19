package v5

import (
	"time"

	v5 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v5"
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

	var sessionTimeoutMS int32
	if sessionTimeoutMS, err = reader.Int32(); err != nil {
		return nil, err
	}
	msg.SessionTimeout = time.Duration(sessionTimeoutMS) * time.Millisecond

	var rebalanceTimeoutMS int32
	if rebalanceTimeoutMS, err = reader.Int32(); err != nil {
		return nil, err
	}
	msg.RebalanceTimeout = time.Duration(rebalanceTimeoutMS) * time.Millisecond

	if msg.MemberID, err = reader.String(); err != nil {
		return nil, err
	}

	if msg.GroupInstanceId, err = reader.NullableString(); err != nil {
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
		groupProtocol := v5.GroupProtocol{}

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
