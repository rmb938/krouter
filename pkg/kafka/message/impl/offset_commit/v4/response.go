package v4

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type OffsetCommitPartitionResponse struct {
	PartitionIndex int32
	ErrCode        errors.KafkaError
}

type OffsetCommitTopicResponse struct {
	Name       string
	Partitions []OffsetCommitPartitionResponse
}

type Response struct {
	ThrottleDuration time.Duration
	Topics           []OffsetCommitTopicResponse
}

func (r *Response) Version() int16 {
	return Version
}
