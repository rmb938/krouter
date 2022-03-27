package v4

type RequestOffsetFetchTopic struct {
	Name             string
	PartitionIndexes []int32
}

type Request struct {
	GroupID string
	Topics  []RequestOffsetFetchTopic
}

func (r *Request) Version() int16 {
	return Version
}
