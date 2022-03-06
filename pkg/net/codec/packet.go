package codec

import (
	"bytes"
	"fmt"
)

type Packet struct {
	ReqHeader *RequestHeader

	Key     *int16
	Version *int16
	Data    []byte
}

func DecodePacket(encodedData []byte) (*Packet, error) {
	rawDecoder := &RawDecoder{Buff: bytes.NewBuffer(encodedData)}

	header := &RequestHeader{}
	err := header.Decode(rawDecoder)
	if err != nil {
		return nil, err
	}

	return &Packet{
		ReqHeader: header,
		Data:      rawDecoder.Buff.Bytes(),
	}, nil
}

func (p *Packet) Encode(correlationId int32) ([]byte, error) {
	header := &ResponseHeader{
		Length:        int32(len(p.Data)),
		CorrelationId: correlationId,
		Version:       *p.Version,
	}

	rawEncoder := &RawEncoder{Buff: bytes.NewBuffer(nil)}
	err := header.Encode(rawEncoder)
	if err != nil {
		return nil, err
	}

	if *p.Key != 3 {
		fmt.Printf("Key: %v Version: %v Length: %v Data: %v\n", *p.Key, header.Version, len(p.Data), p.Data)
	}

	// return header + data
	return append(rawEncoder.ToBytes(), p.Data...), nil
}
