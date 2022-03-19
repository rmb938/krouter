package v5

import "time"

type ListOffsetsPartitionRequest struct {
	PartitionIndex     int32
	CurrentLeaderEpoch int32
	Timestamp          time.Time
}

type ListOffsetsTopicRequest struct {
	Name       string
	Partitions []ListOffsetsPartitionRequest
}

type Request struct {
	ReplicaID      int32
	IsolationLevel int8
	Topics         []ListOffsetsTopicRequest
}

func (r *Request) Version() int16 {
	return Version
}
