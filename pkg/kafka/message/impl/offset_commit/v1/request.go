package v4

import "time"

type OffsetCommitPartitionRequest struct {
	PartitionIndex     int32
	CommittedOffset    int64
	CommittedTimestamp time.Time
	CommittedMetadata  *string
}

type OffsetCommitTopicRequest struct {
	Name       string
	Partitions []OffsetCommitPartitionRequest
}

type Request struct {
	GroupID      string
	GenerationID int32
	MemberID     string
	Topics       []OffsetCommitTopicRequest
}

func (r *Request) Version() int16 {
	return Version
}
