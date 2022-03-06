package v0

import (
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v0.Request{}

	return msg, nil
}
