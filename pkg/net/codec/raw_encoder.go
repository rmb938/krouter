package codec

import (
	"bytes"
	"encoding/binary"
)

type RawEncoder struct {
	Buff *bytes.Buffer
}

func (pe *RawEncoder) Int8(in int8) {
	pe.Buff.WriteByte(byte(in))
}

func (pe *RawEncoder) Int16(in int16) {
	buff := make([]byte, 2)
	binary.BigEndian.PutUint16(buff, uint16(in))
	pe.Buff.Write(buff)
}

func (pe *RawEncoder) Int32(in int32) {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, uint32(in))
	pe.Buff.Write(buff)
}

func (pe *RawEncoder) Int64(in int64) {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, uint64(in))
	pe.Buff.Write(buff)
}

func (pe *RawEncoder) VarInt(in int64) int {
	buff := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buff, in)
	pe.Buff.Write(buff[:n])
	return n
}

func (pe *RawEncoder) UVarInt(in uint64) int {
	buff := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buff, in)
	pe.Buff.Write(buff[:n])
	return n
}

func (pe *RawEncoder) Bool(in bool) {
	if in {
		pe.Int8(1)
	} else {
		pe.Int8(0)
	}
}

func (pe *RawEncoder) ArrayLength(in int) {
	pe.Int32(int32(in))
}

func (pe *RawEncoder) CompactArrayLength(in int) {
	pe.UVarInt(uint64(in + 1))
}

func (pe *RawEncoder) Bytes(in []byte) {
	pe.ArrayLength(len(in))
	pe.Buff.Write(in)
}

func (pe *RawEncoder) CompactBytes(in []byte) {
	pe.CompactArrayLength(len(in))
	pe.Buff.Write(in)
}

func (pe *RawEncoder) NullableBytes(in []byte) {
	if in == nil {
		pe.Int32(-1)
		return
	}

	pe.Bytes(in)
}

func (pe *RawEncoder) CompactNullableBytes(in []byte) {
	if in == nil {
		pe.UVarInt(uint64(0))
		return
	}

	pe.CompactBytes(in)
}

func (pe *RawEncoder) String(in string) {
	pe.Int16(int16(len(in)))
	pe.Buff.Write([]byte(in))
}

func (pe *RawEncoder) NullableString(in *string) {
	if in == nil {
		pe.Int16(-1)
		return
	}

	pe.String(*in)
}

func (pe *RawEncoder) CompactString(in string) {
	pe.CompactArrayLength(len(in))
	pe.Buff.Write([]byte(in))
}

func (pe *RawEncoder) CompactNullableString(in *string) {
	if in == nil {
		pe.UVarInt(0)
		return
	}

	pe.CompactString(*in)
}

func (pe *RawEncoder) ToBytes() []byte {
	return pe.Buff.Bytes()
}
