package v0

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v0.Response)

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

		// resource_type
		builder.Encoder.Int8(result.ResourceType)

		// resource_name
		builder.Encoder.String(result.ResourceName)

		// configs
		builder.Encoder.ArrayLength(len(result.Configs))
		for _, config := range result.Configs {
			// name
			builder.Encoder.String(config.Name)

			// value
			builder.Encoder.NullableString(config.Value)

			// read_only
			builder.Encoder.Bool(config.ReadOnly)

			// is_default
			builder.Encoder.Bool(config.Default)

			// is_sensitive
			builder.Encoder.Bool(config.Sensitive)
		}
	}

	return builder.ToPacket(), nil
}
