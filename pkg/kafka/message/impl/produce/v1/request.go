package v7

import "time"

type PartitionData struct {
	Index   int32
	Records []byte
}

type TopicData struct {
	Name          string
	PartitionData []PartitionData
}

type Request struct {
	ACKs            int16
	TimeoutDuration time.Duration
	TopicData       []TopicData
}

func (r *Request) Version() int16 {
	return Version
}
