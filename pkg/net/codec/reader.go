package codec

type PackerReader struct {
	*RawDecoder
}

func NewPacketReader(packet *Packet) *PackerReader {
	return &PackerReader{
		RawDecoder: &RawDecoder{Buff: packet.Data},
	}
}
