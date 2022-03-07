package records

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var ErrInsufficientData = errors.New("kafka record: insufficient data to decode packet, more bytes expected")
var ErrVarIntOverflow = errors.New("kafka record: varint overflow")

type RecordDecoder struct {
	Buff *bytes.Buffer
}

func (pd *RecordDecoder) Close() error {

	if pd.Buff.Len() > 0 {
		return fmt.Errorf("closed record decoder while there was still data to decode %d bytes left", pd.Buff.Len())
	}

	return nil
}

func (pd *RecordDecoder) Int8() (int8, error) {
	if pd.Buff.Len() < 1 {
		return -1, ErrInsufficientData
	}

	return int8(pd.Buff.Next(1)[0]), nil
}

func (pd *RecordDecoder) Int16() (int16, error) {
	if pd.Buff.Len() < 2 {
		return -1, ErrInsufficientData
	}

	return int16(binary.BigEndian.Uint16(pd.Buff.Next(2))), nil
}

func (pd *RecordDecoder) Int32() (int32, error) {
	if pd.Buff.Len() < 4 {
		return -1, ErrInsufficientData
	}

	return int32(binary.BigEndian.Uint32(pd.Buff.Next(4))), nil
}

func (pd *RecordDecoder) Int64() (int64, error) {
	if pd.Buff.Len() < 8 {
		return -1, ErrInsufficientData
	}

	return int64(binary.BigEndian.Uint64(pd.Buff.Next(8))), nil
}

func (pd *RecordDecoder) VarInt() (int64, error) {
	tmp, n := binary.Varint(pd.Buff.Bytes())
	if n == 0 {
		return -1, ErrInsufficientData
	}
	if n < 0 {
		pd.Buff.Next(0 - n)
		return -1, ErrVarIntOverflow
	}

	pd.Buff.Next(n)
	return tmp, nil
}

func (pd *RecordDecoder) ArrayLength() (int32, error) {
	return pd.Int32()
}

func (pd *RecordDecoder) VarIntBytes() ([]byte, error) {
	length, err := pd.VarInt()
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return nil, nil
	}

	return pd.Buff.Next(int(length)), nil
}

func (pd *RecordDecoder) Length() int {
	return pd.Buff.Len()
}

func (pd *RecordDecoder) RawBytes(size int) []byte {
	return pd.Buff.Next(size)
}
