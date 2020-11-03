package codec

import "bytes"

type PackerReader struct {
	*RawDecoder
}

func NewPacketReader(packet *Packet) *PackerReader {
	return &PackerReader{
		RawDecoder: &RawDecoder{Buff: bytes.NewBuffer(packet.Data)},
	}
}
