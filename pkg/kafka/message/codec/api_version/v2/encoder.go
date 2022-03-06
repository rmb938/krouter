package v2

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v2"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v2.Response)

	builder := codec.NewPacketBuilder(api_version.Key, msg.Version())

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	builder.Encoder.ArrayLength(len(msg.APIKeys))
	for _, apiKey := range msg.APIKeys {
		builder.Encoder.Int16(apiKey.Key)
		builder.Encoder.Int16(apiKey.MinVersion)
		builder.Encoder.Int16(apiKey.MaxVersion)
	}

	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	return builder.ToPacket(), nil
}
