package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

var ErrInsufficientData = errors.New("kafka: insufficient data to decode packet, more bytes expected")
var ErrVarIntOverflow = errors.New("kafka: varint overflow")
var ErrUVarIntOverflow = errors.New("kafka: uvarint overflow")
var ErrInvalidBool = errors.New("kafka: invalid bool")

type RawDecoder struct {
	Buff *bytes.Buffer
}

func (pd *RawDecoder) Close() error {

	if pd.Buff.Len() > 0 {
		return fmt.Errorf("closed raw decoder while there was still data to decode %d bytes left: %v", pd.Buff.Len(), hex.EncodeToString(pd.Buff.Next(pd.Buff.Len())))
	}

	return nil
}

func (pd *RawDecoder) Int8() (int8, error) {
	if pd.Buff.Len() < 1 {
		return -1, ErrInsufficientData
	}

	return int8(pd.Buff.Next(1)[0]), nil
}

func (pd *RawDecoder) Int16() (int16, error) {
	if pd.Buff.Len() < 2 {
		return -1, ErrInsufficientData
	}

	return int16(binary.BigEndian.Uint16(pd.Buff.Next(2))), nil
}

func (pd *RawDecoder) Int32() (int32, error) {
	if pd.Buff.Len() < 4 {
		return -1, ErrInsufficientData
	}

	return int32(binary.BigEndian.Uint32(pd.Buff.Next(4))), nil
}

func (pd *RawDecoder) Int64() (int64, error) {
	if pd.Buff.Len() < 8 {
		return -1, ErrInsufficientData
	}

	return int64(binary.BigEndian.Uint64(pd.Buff.Next(8))), nil
}

func (pd *RawDecoder) VarInt() (int64, error) {
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

func (pd *RawDecoder) UVarInt() (uint64, error) {
	tmp, n := binary.Uvarint(pd.Buff.Bytes())
	if n == 0 {
		return 0, ErrInsufficientData
	}
	if n < 0 {
		pd.Buff.Next(0 - n)
		return 0, ErrUVarIntOverflow
	}

	pd.Buff.Next(n)
	return tmp, nil
}

func (pd *RawDecoder) Bool() (bool, error) {
	tmp, err := pd.Int8()
	if err != nil {
		return false, err
	}

	if tmp == 0 {
		return false, nil
	}

	if tmp != 1 {
		return false, ErrInvalidBool
	}

	return true, nil
}

func (pd *RawDecoder) Bytes() ([]byte, error) {
	length, err := pd.Int32()
	if err != nil {
		return nil, err
	}

	return pd.Buff.Next(int(length)), nil
}

func (pd *RawDecoder) CompactBytes() ([]byte, error) {
	length, err := pd.UVarInt()
	if err != nil {
		return nil, err
	}

	return pd.Buff.Next(int(length - 1)), nil
}

func (pd *RawDecoder) NullableBytes() ([]byte, error) {
	length, err := pd.Int32()
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return nil, nil
	}

	return pd.Buff.Next(int(length)), nil
}

func (pd *RawDecoder) CompactNullableBytes() ([]byte, error) {
	length, err := pd.UVarInt()
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}

	return pd.Buff.Next(int(length - 1)), nil
}

func (pd *RawDecoder) String() (string, error) {
	length, err := pd.Int16()
	if err != nil {
		return "", err
	}

	return string(pd.Buff.Next(int(length))), nil
}

func (pd *RawDecoder) NullableString() (*string, error) {
	length, err := pd.Int16()
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return nil, nil
	}

	tmp := string(pd.Buff.Next(int(length)))
	return &tmp, nil
}

func (pd *RawDecoder) CompactString() (string, error) {
	length, err := pd.CompactArrayLength()
	if err != nil {
		return "", err
	}

	return string(pd.Buff.Next(int(length - 1))), nil
}

func (pd *RawDecoder) CompactNullableString() (*string, error) {
	length, err := pd.UVarInt()
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}

	tmp := string(pd.Buff.Next(int(length)))
	return &tmp, nil
}

func (pd *RawDecoder) ArrayLength() (int32, error) {
	return pd.Int32()
}

func (pd *RawDecoder) NullableArrayLength() (int32, error) {
	length, err := pd.Int32()
	if err != nil {
		return 0, err
	}

	if length < 0 {
		return 0, nil
	}

	return length, nil
}

func (pd *RawDecoder) CompactArrayLength() (uint64, error) {
	return pd.UVarInt()
}

func (pd *RawDecoder) Records() ([]byte, error) {
	return pd.NullableBytes()
}

func (pd *RawDecoder) CompactRecords() ([]byte, error) {
	return pd.CompactNullableBytes()
}
