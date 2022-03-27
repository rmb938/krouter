package v1

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_fetch/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v1.Response)

	builder := codec.NewPacketBuilder(metadata.Key, msg.Version())

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

	return builder.ToPacket(), nil
}
