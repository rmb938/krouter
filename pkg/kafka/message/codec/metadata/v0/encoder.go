package v0

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/metadata"
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v0.Response)

	builder := codec.NewPacketBuilder(metadata.Key, msg.Version())

	// brokers
	builder.Encoder.ArrayLength(len(msg.Brokers))
	for _, broker := range msg.Brokers {
		builder.Encoder.Int32(broker.ID)    // node id
		builder.Encoder.String(broker.Host) // host
		builder.Encoder.Int32(broker.Port)  // port
	}

	// topics
	builder.Encoder.ArrayLength(len(msg.Topics))
	for _, topic := range msg.Topics {
		builder.Encoder.Int16(int16(topic.ErrCode)) // error code
		builder.Encoder.String(topic.Name)          // name

		// partitions
		builder.Encoder.ArrayLength(len(topic.Partitions))
		for _, partition := range topic.Partitions {
			builder.Encoder.Int16(int16(partition.ErrCode)) // error code
			builder.Encoder.Int32(partition.Index)          // partition index
			builder.Encoder.Int32(partition.LeaderID)       // leader id

			// replica nodes
			builder.Encoder.ArrayLength(len(partition.ReplicaNodes))
			for _, replicaNode := range partition.ReplicaNodes {
				builder.Encoder.Int32(replicaNode)
			}

			// isr nodes
			builder.Encoder.ArrayLength(len(partition.ISRNodes))
			for _, isrNode := range partition.ISRNodes {
				builder.Encoder.Int32(isrNode)
			}
		}
	}

	return builder.ToPacket(), nil
}
