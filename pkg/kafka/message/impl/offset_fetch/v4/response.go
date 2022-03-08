package v4

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type ResponseOffsetFetchTopicPartition struct {
	PartitionIndex  int32
	CommittedOffset int64
	Metadata        *string
	ErrCode         errors.KafkaError
}

type ResponseOffsetFetchTopic struct {
	Name       string
	Partitions []ResponseOffsetFetchTopicPartition
}

type Response struct {
	ThrottleDuration time.Duration
	Topics           []ResponseOffsetFetchTopic
	ErrCode          errors.KafkaError
}

func (r *Response) Version() int16 {
	return Version
}
