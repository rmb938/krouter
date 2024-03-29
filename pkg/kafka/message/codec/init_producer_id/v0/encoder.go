package v0

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v0.Response)

	builder := codec.NewPacketBuilder(init_producer_id.Key, msg.Version())

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	// producer_id
	builder.Encoder.Int64(msg.ProducerID)

	// producer_epoch
	builder.Encoder.Int16(msg.ProducerEpoch)

	return builder.ToPacket(), nil
}
