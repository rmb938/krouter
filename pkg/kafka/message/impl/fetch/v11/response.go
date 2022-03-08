package v11

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type FetchAbortedTransaction struct {
	ProducerID  int64
	FirstOffset int64
}

type FetchPartitionResponse struct {
	PartitionIndex       int32
	ErrCode              errors.KafkaError
	HighWaterMark        int64
	LastStableOffset     int64
	LogStartOffset       int64
	AbortedTransactions  []FetchAbortedTransaction
	PreferredReadReplica int32
	Records              []byte
}

type FetchTopicResponse struct {
	Topic      string
	Partitions []FetchPartitionResponse
}

type Response struct {
	ThrottleDuration time.Duration
	ErrCode          errors.KafkaError
	SessionID        int32
	Responses        []FetchTopicResponse
}

func (r *Response) Version() int16 {
	return Version
}
