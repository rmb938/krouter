package records

import (
	"bytes"
	"encoding/binary"
)

type RecordEncoder struct {
	Buff *bytes.Buffer
}

func (pe *RecordEncoder) Int8(in int8) {
	pe.Buff.WriteByte(byte(in))
}

func (pe *RecordEncoder) Int16(in int16) {
	buff := make([]byte, 2)
	binary.BigEndian.PutUint16(buff, uint16(in))
	pe.Buff.Write(buff)
}

func (pe *RecordEncoder) Int32(in int32) {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, uint32(in))
	pe.Buff.Write(buff)
}

func (pe *RecordEncoder) Int64(in int64) {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, uint64(in))
	pe.Buff.Write(buff)
}

func (pe *RecordEncoder) VarInt(in int64) int {
	buff := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buff, in)
	pe.Buff.Write(buff[:n])
	return n
}

func (pe *RecordEncoder) ArrayLength(in int) {
	pe.Int32(int32(in))
}

func (pe *RecordEncoder) VarIntBytes(in []byte) {
	if len(in) == 0 {
		pe.VarInt(-1)
		return
	}
	pe.VarInt(int64(len(in)))
	pe.Buff.Write(in)
}

func (pe *RecordEncoder) RawBytes(in []byte) {
	pe.Buff.Write(in)
}

func (pe *RecordEncoder) ToBytes() []byte {
	return pe.Buff.Bytes()
}
