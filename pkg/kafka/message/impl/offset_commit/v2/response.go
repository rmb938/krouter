package v4

import (
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
	Topics []OffsetCommitTopicResponse
}

func (r *Response) Version() int16 {
	return Version
}
