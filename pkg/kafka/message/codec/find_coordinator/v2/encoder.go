package v2

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator/v2"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v2.Response)

	builder := codec.NewPacketBuilder(find_coordinator.Key, msg.Version())

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	// error_message
	builder.Encoder.NullableString(msg.ErrMessage)

	// node_id
	builder.Encoder.Int32(msg.NodeID)

	// host
	builder.Encoder.String(msg.Host)

	// port
	builder.Encoder.Int32(msg.Port)

	return builder.ToPacket(), nil
}
