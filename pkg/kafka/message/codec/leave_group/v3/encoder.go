package v3

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	v3 "github.com/rmb938/krouter/pkg/kafka/message/impl/leave_group/v3"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v3.Response)

	builder := codec.NewPacketBuilder(find_coordinator.Key, msg.Version())

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	// members
	builder.Encoder.ArrayLength(len(msg.Members))
	for _, member := range msg.Members {
		// member_id
		builder.Encoder.String(member.MemberID)

		// group_instance_id
		builder.Encoder.NullableString(member.GroupInstanceID)

		// error_code
		builder.Encoder.Int16(int16(member.ErrCode))
	}

	return builder.ToPacket(), nil
}
