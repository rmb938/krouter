package v1

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/create_acls/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v1.Response)

	builder := codec.NewPacketBuilder(api_version.Key, msg.Version())

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	// results
	builder.Encoder.ArrayLength(len(msg.Results))
	for _, result := range msg.Results {
		// error_code
		builder.Encoder.Int16(int16(result.ErrCode))

		// error_message
		builder.Encoder.NullableString(result.ErrMessage)
	}

	return builder.ToPacket(), nil
}
