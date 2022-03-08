package v4

import "time"

type OffsetCommitPartitionRequest struct {
	PartitionIndex    int32
	CommittedOffset   int64
	CommittedMetadata *string
}

type OffsetCommitTopicRequest struct {
	Name       string
	Partitions []OffsetCommitPartitionRequest
}

type Request struct {
	GroupID       string
	GenerationID  int32
	MemberID      string
	RetentionTime time.Duration
	Topics        []OffsetCommitTopicRequest
}

func (r *Request) Version() int16 {
	return Version
}
