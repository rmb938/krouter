package v0

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/heartbeat/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v0.Response)

	builder := codec.NewPacketBuilder(find_coordinator.Key, msg.Version())

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	return builder.ToPacket(), nil
}
