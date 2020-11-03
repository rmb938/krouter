package message

import (
	"github.com/rmb938/krouter/pkg/net/codec"
)

type Message interface {
	Version() int16
}

type Decoder interface {
	Decode(reader *codec.PackerReader) (Message, error)
}

type Encoder interface {
	Encode(message Message) (*codec.Packet, error)
}
