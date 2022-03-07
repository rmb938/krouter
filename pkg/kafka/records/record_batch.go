package records

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"hash/crc32"
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

	length, err := decoder.Int32()
	if err != nil {
		return nil, err
	}
	_ = length // TODO: do we need to do anything with this?

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

	crc, err := decoder.Int32()
	if err != nil {
		return nil, err
	}
	_ = crc
	// TODO: validate CRC somehow

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

		length, err := recordDecoder.VarInt()
		if err != nil {
			return nil, err
		}
		_ = length // TODO: do we need to do anything with this?

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

		record.Key, err = recordDecoder.VarIntBytes()
		if err != nil {
			return nil, err
		}

		record.Value, err = recordDecoder.VarIntBytes()
		if err != nil {
			return nil, err
		}

		headerCount, err := recordDecoder.VarInt()
		if err != nil {
			return nil, err
		}

		for i := int64(0); i < headerCount; i++ {
			header := RecordHeader{}

			header.Key, err = recordDecoder.VarIntBytes()
			if err != nil {
				return nil, err
			}

			header.Value, err = recordDecoder.VarIntBytes()
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

func (rb *RecordBatch) Encode() ([]byte, error) {
	encoder := RecordEncoder{Buff: bytes.NewBuffer([]byte{})}

	encoder.Int16(rb.Attributes)

	encoder.Int32(rb.LastOffsetDelta)

	encoder.Int64(rb.FirstTimestamp.UnixMilli())

	encoder.Int64(rb.MaxTimestamp.UnixMilli())

	encoder.Int64(rb.ProducerID)

	encoder.Int16(rb.ProducerEpoch)

	encoder.Int32(rb.BaseSequence)

	encoder.ArrayLength(len(rb.Records))

	recordsEncoder := RecordEncoder{Buff: bytes.NewBuffer([]byte{})}

	for _, record := range rb.Records {
		recordEncoder := RecordEncoder{Buff: bytes.NewBuffer([]byte{})}

		recordEncoder.Int8(record.Attributes)

		recordEncoder.VarInt(int64(record.TimeStampDelta))

		recordEncoder.VarInt(record.OffsetDelta)

		recordEncoder.VarIntBytes(record.Key)

		recordEncoder.VarIntBytes(record.Value)

		recordEncoder.VarInt(int64(len(record.Headers)))
		for _, header := range record.Headers {
			recordEncoder.VarIntBytes(header.Key)
			recordEncoder.VarIntBytes(header.Value)
		}

		recordHeaderEncoder := RecordEncoder{Buff: bytes.NewBuffer([]byte{})}
		recordBytes := recordEncoder.ToBytes()
		recordHeaderEncoder.VarInt(int64(len(recordBytes)))
		recordHeaderEncoder.RawBytes(recordBytes)

		recordsEncoder.RawBytes(recordHeaderEncoder.ToBytes())
	}

	compressedRecordBytes, err := rb.compress(recordsEncoder.ToBytes())
	if err != nil {
		return nil, err
	}

	encoder.RawBytes(compressedRecordBytes)

	header2Encoder := RecordEncoder{Buff: bytes.NewBuffer([]byte{})}
	header2Encoder.Int32(rb.PartitionLeaderEpoch)
	header2Encoder.Int8(2) // hardcode magic to 2

	encoderBytes := encoder.ToBytes()
	header2Encoder.Int32(int32(crc32.Checksum(encoderBytes, crc32.MakeTable(crc32.Castagnoli))))
	header2Encoder.RawBytes(encoderBytes)

	header1Encoder := RecordEncoder{Buff: bytes.NewBuffer([]byte{})}
	header1Encoder.Int64(rb.BaseOffset)

	header2Bytes := header2Encoder.ToBytes()

	header1Encoder.Int32(int32(len(header2Bytes)))
	header1Encoder.RawBytes(header2Bytes)

	return encoder.ToBytes(), nil
}

func (rb *RecordBatch) compress(data []byte) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})

	switch rb.GetCodec() {
	case CompressionNone:
		return data, nil
	case CompressionGZIP:
		writer := gzip.NewWriter(bufio.NewWriter(bufio.NewWriter(output)))
		defer writer.Close()

		_, err := writer.Write(data)
		if err != nil {
			return nil, err
		}
	case CompressionSnappy:
		return snappy.Encode(data), nil
	case CompressionLZ4:
		writer := lz4.NewWriter(bufio.NewWriter(output))
		defer writer.Close()

		_, err := writer.Write(data)
		if err != nil {
			return nil, err
		}
	case CompressionZSTD:
		zstdEnc, err := zstd.NewWriter(nil, zstd.WithZeroFrames(true), zstd.WithEncoderLevel(zstd.SpeedDefault))
		if err != nil {
			return nil, err
		}
		defer zstdEnc.Close()

		_, err = zstdEnc.Write(data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid compression specified: %d", rb.GetCodec())
	}

	return output.Bytes(), nil
}
