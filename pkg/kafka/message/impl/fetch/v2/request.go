package v11

import "time"

type FetchPartitionRequest struct {
	Partition         int32
	FetchOffset       int64
	PartitionMaxBytes int32
}

type FetchTopicRequest struct {
	Name       string
	Partitions []FetchPartitionRequest
}

type Request struct {
	MaxWait  time.Duration
	MinBytes int32
	Topics   []FetchTopicRequest
}

func (r *Request) Version() int16 {
	return Version
}
