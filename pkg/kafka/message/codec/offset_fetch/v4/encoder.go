package v4

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	v4 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v4"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v4.Response)

	builder := codec.NewPacketBuilder(metadata.Key, msg.Version())

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration / time.Millisecond))

	// topics
	builder.Encoder.ArrayLength(len(msg.Topics))
	for _, topic := range msg.Topics {
		// name
		builder.Encoder.String(topic.Name)

		// partitions
		builder.Encoder.ArrayLength(len(topic.Partitions))
		for _, partition := range topic.Partitions {
			// partition_index
			builder.Encoder.Int32(partition.PartitionIndex)

			// committed_offset
			builder.Encoder.Int64(partition.CommittedOffset)

			// metadata
			builder.Encoder.NullableString(partition.Metadata)

			// error_code
			builder.Encoder.Int16(int16(partition.ErrCode))
		}
	}

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	return builder.ToPacket(), nil
}
