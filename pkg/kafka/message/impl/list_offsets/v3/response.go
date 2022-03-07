package v3

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type ListOffsetsPartitionResponse struct {
	PartitionIndex int32
	ErrCode        errors.KafkaError
	Timestamp      time.Time
	Offset         int64
}

type ListOffsetsTopicResponse struct {
	Name       string
	Partitions []ListOffsetsPartitionResponse
}

type Response struct {
	ThrottleDuration time.Duration
	Topics           []ListOffsetsTopicResponse
}

func (r *Response) Version() int16 {
	return Version
}
