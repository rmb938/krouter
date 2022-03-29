package v1

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_acls/v1"
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

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	// error_message
	builder.Encoder.NullableString(msg.ErrMessage)

	// resources
	builder.Encoder.ArrayLength(len(msg.Resources))
	for _, resource := range msg.Resources {
		// resource_type
		builder.Encoder.Int8(int8(resource.ResourceType))

		// resource_name
		builder.Encoder.String(resource.ResourceName)

		// pattern_type
		builder.Encoder.Int8(int8(resource.PatternType))

		// acls
		builder.Encoder.ArrayLength(len(resource.ACLs))
		for _, acl := range resource.ACLs {
			// principal
			builder.Encoder.String(acl.Principal)

			// host
			builder.Encoder.String(acl.Host)

			// operation
			builder.Encoder.Int8(int8(acl.Operation))

			// permission_type
			builder.Encoder.Int8(int8(acl.PermissionType))
		}
	}

	return builder.ToPacket(), nil
}
