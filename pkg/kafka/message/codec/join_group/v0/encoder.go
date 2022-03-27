package v0

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/join_group/v0"
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

	// generation_id
	builder.Encoder.Int32(msg.GenerationID)

	// protocol_name
	builder.Encoder.String(msg.ProtocolName)

	// leader
	builder.Encoder.String(msg.Leader)

	// member_id
	builder.Encoder.String(msg.MemberID)

	// members
	builder.Encoder.ArrayLength(len(msg.Members))
	for _, member := range msg.Members {
		// member_id
		builder.Encoder.String(member.MemberID)

		// metadata
		builder.Encoder.Bytes(member.Metadata)
	}

	return builder.ToPacket(), nil
}
