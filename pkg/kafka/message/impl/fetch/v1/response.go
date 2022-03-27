package v11

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type FetchPartitionResponse struct {
	PartitionIndex int32
	ErrCode        errors.KafkaError
	HighWaterMark  int64
	Records        []byte
}

type FetchTopicResponse struct {
	Topic      string
	Partitions []FetchPartitionResponse
}

type Response struct {
	ThrottleDuration time.Duration
	Responses        []FetchTopicResponse
}

func (r *Response) Version() int16 {
	return Version
}
