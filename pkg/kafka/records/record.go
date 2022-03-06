package records

import "time"

type RecordHeader struct {
	Key   []byte
	Value []byte
}

type Record struct {
	Attributes     int8
	TimeStampDelta time.Duration
	OffsetDelta    int64
	Key            []byte
	Value          []byte
	Headers        []RecordHeader
}
