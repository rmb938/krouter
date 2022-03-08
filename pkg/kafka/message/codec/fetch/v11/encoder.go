package v11

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/api_version"
	v11 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v11"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Encoder struct {
}

func (e *Encoder) Encode(message message.Message) (*codec.Packet, error) {
	msg := message.(*v11.Response)

	builder := codec.NewPacketBuilder(api_version.Key, msg.Version())

	// throttle_time_ms
	builder.Encoder.Int32(int32(msg.ThrottleDuration.Milliseconds()))

	// error_code
	builder.Encoder.Int16(int16(msg.ErrCode))

	// session_id
	builder.Encoder.Int32(msg.SessionID)

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

			// last_stable_offset
			builder.Encoder.Int64(partition.LastStableOffset)

			// log_start_offset
			builder.Encoder.Int64(partition.LogStartOffset)

			// aborted_transactions
			builder.Encoder.ArrayLength(len(partition.AbortedTransactions))
			for _, abortedTransaction := range partition.AbortedTransactions {
				// producer_id
				builder.Encoder.Int64(abortedTransaction.ProducerID)

				// first_offset
				builder.Encoder.Int64(abortedTransaction.FirstOffset)
			}

			// preferred_read_replica
			builder.Encoder.Int32(partition.PreferredReadReplica)

			// records
			builder.Encoder.Records(partition.Records)
		}
	}

	return builder.ToPacket(), nil
}
