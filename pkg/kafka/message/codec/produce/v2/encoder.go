package v2

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/produce"
	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v2"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v2.Response)

	builder := codec.NewPacketBuilder(produce.Key, msg.Version())

	// responses
	builder.Encoder.ArrayLength(len(msg.Responses))
	for _, response := range msg.Responses {
		// name
		builder.Encoder.String(response.Name)

		// partition_responses
		builder.Encoder.ArrayLength(len(response.PartitionResponses))

		for _, partitionResponse := range response.PartitionResponses {
			// index
			builder.Encoder.Int32(partitionResponse.Index)

			// error_code
			builder.Encoder.Int16(int16(partitionResponse.ErrCode))

			// base_offset
			builder.Encoder.Int64(partitionResponse.BaseOffset)

			// log_append_time_ms
			builder.Encoder.Int64(partitionResponse.LogAppendTime.Unix())
		}
	}

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	return builder.ToPacket(), nil
}
