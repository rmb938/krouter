package v5

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type PartitionResponse struct {
	Index          int32
	ErrCode        errors.KafkaError
	BaseOffset     int64
	LogAppendTime  time.Time
	LogStartOffset int64
}

type ProduceResponse struct {
	Name               string
	PartitionResponses []PartitionResponse
}

type Response struct {
	Responses        []ProduceResponse
	ThrottleDuration time.Duration
}

func (r *Response) Version() int16 {
	return Version
}
