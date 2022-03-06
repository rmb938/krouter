package records

import (
	"bytes"
	"time"
)

type RecordBatch struct {
	BaseOffset           int64
	PartitionLeaderEpoch int32
	Magic                int8
	CRC                  int32
	Attributes           int16
	LastOffsetDelta      int32
	FirstTimestamp       time.Time
	MaxTimestamp         time.Time
	ProducerID           int64
	ProducerEpoch        int16
	BaseSequence         int32
	Records              []Record
}

func ParseRecordBatch(recordBatchBytes []byte) (*RecordBatch, error) {
	rb := &RecordBatch{}

	decoder := &RecordDecoder{Buff: bytes.NewBuffer(recordBatchBytes)}

	var err error
	rb.BaseOffset, err = decoder.Int64()
	if err != nil {
		return nil, err
	}

	_, err = decoder.Int32()
	if err != nil {
		return nil, err
	}

	rb.PartitionLeaderEpoch, err = decoder.Int32()
	if err != nil {
		return nil, err
	}

	rb.Magic, err = decoder.Int8()
	if err != nil {
		return nil, err
	}

	_, err = decoder.Int32()
	if err != nil {
		return nil, err
	}

	rb.Attributes, err = decoder.Int16()
	if err != nil {
		return nil, err
	}

	rb.LastOffsetDelta, err = decoder.Int32()
	if err != nil {
		return nil, err
	}

	firstTimestampMS, err := decoder.Int64()
	if err != nil {
		return nil, err
	}
	rb.FirstTimestamp = time.UnixMilli(firstTimestampMS)

	maxTimestamp, err := decoder.Int64()
	if err != nil {
		return nil, err
	}
	rb.MaxTimestamp = time.UnixMilli(maxTimestamp)

	rb.ProducerID, err = decoder.Int64()
	if err != nil {
		return nil, err
	}

	rb.ProducerEpoch, err = decoder.Int16()
	if err != nil {
		return nil, err
	}

	rb.BaseSequence, err = decoder.Int32()
	if err != nil {
		return nil, err
	}

	recordCount, err := decoder.ArrayLength()
	if err != nil {
		return nil, err
	}

	for i := int32(0); i < recordCount; i++ {
		record := &Record{}

		_, err := decoder.VarInt()
		if err != nil {
			return nil, err
		}

		record.Attributes, err = decoder.Int8()
		if err != nil {
			return nil, err
		}

		record.TimeStampDelta, err = decoder.VarInt()
		if err != nil {
			return nil, err
		}

		record.OffsetDelta, err = decoder.VarInt()
		if err != nil {
			return nil, err
		}

		record.Key, err = decoder.VarinBytes()
		if err != nil {
			return nil, err
		}

		record.Value, err = decoder.VarinBytes()
		if err != nil {
			return nil, err
		}

		headerCount, err := decoder.VarInt()
		if err != nil {
			return nil, err
		}

		for i := int64(0); i < headerCount; i++ {
			header := RecordHeader{}

			header.Key, err = decoder.VarinBytes()
			if err != nil {
				return nil, err
			}

			header.Value, err = decoder.VarinBytes()
			if err != nil {
				return nil, err
			}

			record.Headers = append(record.Headers, header)
		}

		rb.Records = append(rb.Records, *record)
	}

	err = decoder.Close()
	if err != nil {
		return nil, err
	}

	return rb, nil
}
