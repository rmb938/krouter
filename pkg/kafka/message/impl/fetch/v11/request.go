package v11

import "time"

type FetchPartitionRequest struct {
	Partition          int32
	CurrentLeaderEpoch int32
	FetchOffset        int64
	LogStartOffset     int64
	PartitionMaxBytes  int32
}

type FetchTopicRequest struct {
	Name       string
	Partitions []FetchPartitionRequest
}

type FetchForgottenTopicData struct {
	Name       string
	Partitions []int32
}

type Request struct {
	MaxWait             time.Duration
	MinBytes            int32
	MaxBytes            int32
	IsolationLevel      int8
	SessionID           int32
	SessionEpoch        int32
	Topics              []FetchTopicRequest
	ForgottenTopicsData []FetchForgottenTopicData
	RackID              string
}

func (r *Request) Version() int16 {
	return Version
}
