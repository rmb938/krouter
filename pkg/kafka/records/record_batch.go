package records

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"time"

	snappy "github.com/eapache/go-xerial-snappy"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
)

const (
	// CompressionNone no compression
	CompressionNone CompressionCodec = iota
	// CompressionGZIP compression using GZIP
	CompressionGZIP
	// CompressionSnappy compression using snappy
	CompressionSnappy
	// CompressionLZ4 compression using LZ4
	CompressionLZ4
	// CompressionZSTD compression using ZSTD
	CompressionZSTD

	isTransactionalMask = 0x10
	controlMask         = 0x20

	// The lowest 3 bits contain the compression codec used for the message
	compressionCodecMask int8 = 0x07

	// Bit 3 set for "LogAppend" timestamps
	timestampTypeMask = 0x08
)

var UnsupportedMagic = fmt.Errorf("record batch magic number is unsupported")

type CompressionCodec int8

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
	Records              []*Record
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
	if rb.Magic != int8(2) {
		return nil, UnsupportedMagic
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

	recordBytes := decoder.RawBytes(decoder.Length())

	// need to decompress to read it into a record list
	//  sadly we can't easily just give sarama the raw bytes,
	//  so we need to decompress into records, convert to sarama records
	//  then sarama will recompress
	//  This is a waste of CPU time, but I don't know any other good way
	decompressedBytes, err := rb.decompress(recordBytes)
	if err != nil {
		return nil, err
	}

	recordDecoder := &RecordDecoder{Buff: bytes.NewBuffer(decompressedBytes)}

	for i := int32(0); i < recordCount; i++ {
		record := &Record{}

		_, err := recordDecoder.VarInt()
		if err != nil {
			return nil, err
		}

		record.Attributes, err = recordDecoder.Int8()
		if err != nil {
			return nil, err
		}

		timestampDelta, err := recordDecoder.VarInt()
		if err != nil {
			return nil, err
		}
		record.TimeStampDelta = time.Duration(timestampDelta)

		record.OffsetDelta, err = recordDecoder.VarInt()
		if err != nil {
			return nil, err
		}

		record.Key, err = recordDecoder.VarinBytes()
		if err != nil {
			return nil, err
		}

		record.Value, err = recordDecoder.VarinBytes()
		if err != nil {
			return nil, err
		}

		headerCount, err := recordDecoder.VarInt()
		if err != nil {
			return nil, err
		}

		for i := int64(0); i < headerCount; i++ {
			header := RecordHeader{}

			header.Key, err = recordDecoder.VarinBytes()
			if err != nil {
				return nil, err
			}

			header.Value, err = recordDecoder.VarinBytes()
			if err != nil {
				return nil, err
			}

			record.Headers = append(record.Headers, header)
		}

		rb.Records = append(rb.Records, record)
	}

	err = recordDecoder.Close()
	if err != nil {
		return nil, err
	}

	err = decoder.Close()
	if err != nil {
		return nil, err
	}

	return rb, nil
}

func (rb *RecordBatch) GetCodec() CompressionCodec {
	return CompressionCodec(int8(rb.Attributes) & compressionCodecMask)
}

func (rb *RecordBatch) IsControl() bool {
	return rb.Attributes&controlMask == controlMask
}

func (rb *RecordBatch) IsLogAppendTime() bool {
	return rb.Attributes&timestampTypeMask == timestampTypeMask
}

func (rb *RecordBatch) IsTransactional() bool {
	return rb.Attributes&isTransactionalMask == isTransactionalMask
}

func (rb *RecordBatch) decompress(data []byte) ([]byte, error) {
	switch rb.GetCodec() {
	case CompressionNone:
		return data, nil
	case CompressionGZIP:
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		return io.ReadAll(reader)
	case CompressionSnappy:
		return snappy.Decode(data)
	case CompressionLZ4:
		reader := lz4.NewReader(bytes.NewReader(data))

		return io.ReadAll(reader)
	case CompressionZSTD:
		zstdDec, err := zstd.NewReader(nil)
		if err != nil {
			return nil, err
		}

		return zstdDec.DecodeAll(data, nil)
	default:
		return nil, fmt.Errorf("invalid compression specified: %d", rb.GetCodec())
	}
}
