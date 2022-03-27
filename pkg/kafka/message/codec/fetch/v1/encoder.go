package v1

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v1"
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

	// responses
	builder.Encoder.ArrayLength(len(msg.Responses))
	for _, resp := range msg.Responses {
		// topic
		builder.Encoder.String(resp.Topic)

		// partitions
		builder.Encoder.ArrayLength(len(resp.Partitions))
		for _, partition := range resp.Partitions {
			// partition_index
			builder.Encoder.Int32(partition.PartitionIndex)

			// error_code
			builder.Encoder.Int16(int16(partition.ErrCode))

			// high_watermark
			builder.Encoder.Int64(partition.HighWaterMark)

			// records
			builder.Encoder.Records(partition.Records)
		}
	}

	return builder.ToPacket(), nil
}
