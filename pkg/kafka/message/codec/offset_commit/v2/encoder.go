package v2

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/find_coordinator"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/offset_commit/v2"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v2.Response)

	builder := codec.NewPacketBuilder(find_coordinator.Key, msg.Version())

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

			// error_code
			builder.Encoder.Int16(int16(partition.ErrCode))
		}
	}

	return builder.ToPacket(), nil
}
