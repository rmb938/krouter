package v0

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_groups/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v0.Response)

	builder := codec.NewPacketBuilder(api_version.Key, msg.Version())

	// groups
	builder.Encoder.ArrayLength(len(msg.Groups))
	for _, group := range msg.Groups {
		// error code
		builder.Encoder.Int16(int16(group.ErrCode))

		// group_id
		builder.Encoder.String(group.GroupID)

		// group_state
		builder.Encoder.String(group.GroupState)

		// protocol_type
		builder.Encoder.String(group.ProtocolType)

		// protocol_data
		builder.Encoder.String(group.ProtocolData)

		// members
		builder.Encoder.ArrayLength(len(group.Members))
		for _, member := range group.Members {
			// member_id
			builder.Encoder.String(member.MemberID)

			// client_id
			builder.Encoder.String(member.ClientID)

			// client_host
			builder.Encoder.String(member.ClientHost)

			// member_metadata
			builder.Encoder.Bytes(member.MemberMetadata)

			// member_assignment
			builder.Encoder.Bytes(member.MemberAssignment)
		}
	}

	return builder.ToPacket(), nil
}
