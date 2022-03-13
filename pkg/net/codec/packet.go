package codec

import (
	"bytes"
)

type Packet struct {
	ReqHeader *RequestHeader

	Key     *int16
	Version *int16
	Data    *bytes.Buffer
}

func DecodePacket(encodedData []byte) (*Packet, error) {
	buff := bytes.NewBuffer(encodedData)
	rawDecoder := &RawDecoder{Buff: buff}

	header := &RequestHeader{}
	err := header.Decode(rawDecoder)
	if err != nil {
		return nil, err
	}

	return &Packet{
		ReqHeader: header,
		Data:      buff,
	}, nil
}

func (p *Packet) Encode(correlationId int32) ([]byte, error) {
	header := &ResponseHeader{
		Length:        int32(p.Data.Len()),
		CorrelationId: correlationId,
		Version:       *p.Version,
	}

	rawEncoder := &RawEncoder{Buff: bytes.NewBuffer(make([]byte, 0, header.Length+8))}
	err := header.Encode(rawEncoder)
	if err != nil {
		return nil, err
	}
	rawEncoder.RawBytes(p.Data.Bytes())

	// return header + data
	return rawEncoder.ToBytes(), nil
}
