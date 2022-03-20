package codec

import "bytes"

type PacketBuilder struct {
	key     int16
	version int16
	Encoder RawEncoder
}

func NewPacketBuilder(key int16, version int16) *PacketBuilder {
	return &PacketBuilder{
		key:     key,
		version: version,
		Encoder: RawEncoder{Buff: bytes.NewBuffer(nil)},
	}
}

func (builder *PacketBuilder) ToPacket() *Packet {
	return &Packet{
		Key:     &builder.key,
		Version: &builder.version,
		Data:    builder.Encoder.Buff,
	}
}
