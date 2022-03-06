package records

type RecordHeader struct {
	Key   []byte
	Value []byte
}

type Record struct {
	Attributes     int8
	TimeStampDelta int64
	OffsetDelta    int64
	Key            []byte
	Value          []byte
	Headers        []RecordHeader
}
