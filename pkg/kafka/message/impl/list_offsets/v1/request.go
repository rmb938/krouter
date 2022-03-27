package v5

import "time"

type ListOffsetsPartitionRequest struct {
	PartitionIndex int32
	Timestamp      time.Time
}

type ListOffsetsTopicRequest struct {
	Name       string
	Partitions []ListOffsetsPartitionRequest
}

type Request struct {
	ReplicaID int32
	Topics    []ListOffsetsTopicRequest
}

func (r *Request) Version() int16 {
	return Version
}
