package v3

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/sync_group/v3"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v3.Response)

	builder := codec.NewPacketBuilder(produce.Key, msg.Version())

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	// assignment
	builder.Encoder.Bytes(msg.Assignment)

	return builder.ToPacket(), nil
}
