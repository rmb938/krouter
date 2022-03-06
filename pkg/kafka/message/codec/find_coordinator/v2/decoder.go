package v2

import (
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v2.Request{}

	var err error
	if msg.Key, err = reader.String(); err != nil {
		return nil, err
	}

	if msg.KeyType, err = reader.Int8(); err != nil {
		return nil, err
	}

	return msg, nil
}
