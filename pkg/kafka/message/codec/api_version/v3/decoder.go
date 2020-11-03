package v3

import (
	"fmt"

	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v3"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v3.Request{}

	fmt.Printf("%v\n", reader.Buff.Bytes())
	fmt.Printf("%v\n", string(reader.Buff.Bytes()))

	var err error
	if msg.ClientSoftwareName, err = reader.CompactString(); err != nil {
		return nil, err
	}
	if msg.ClientSoftwareVersion, err = reader.CompactString(); err != nil {
		return nil, err
	}

	// tagged fields
	if _, err := reader.UVarInt(); err != nil {
		return nil, err
	}

	return msg, nil
}
