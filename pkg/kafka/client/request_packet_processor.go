package client

import (
	"encoding/binary"
	"fmt"

	"github.com/go-logr/logr"
	netCodec "github.com/rmb938/krouter/pkg/net/codec"
)

type RequestPacketProcessor struct {
	Log logr.Logger
}

func (pp *RequestPacketProcessor) ProcessPacket(client *Client) (*netCodec.Packet, error) {
	log := pp.Log.WithValues("from-address", client.RemoteAddr().String())

	log.V(1).Info("Reading Packet Length")
	packetLength, err := pp.readPacketLength(client)
	if err != nil {
		return nil, fmt.Errorf("error reading packet length: %w", err)
	}

	log.V(1).Info("Reading Packet Body", "length", packetLength)
	inPacketEncoded := make([]byte, packetLength)
	if _, err := client.ReadFull(inPacketEncoded); err != nil {
		return nil, fmt.Errorf("error reading body of packet: %w", err)
	}

	inPacket, err := netCodec.DecodePacket(inPacketEncoded)
	if err != nil {
		return nil, fmt.Errorf("error decoding packet: %w", err)
	}
	return inPacket, nil
}

func (pp *RequestPacketProcessor) readPacketLength(client *Client) (int64, error) {
	packetLengthBytes := make([]byte, 4)
	if _, err := client.ReadFull(packetLengthBytes); err != nil {
		return 0, fmt.Errorf("error reading length of packet: %w", err)
	}

	packetLength := binary.BigEndian.Uint32(packetLengthBytes)
	if packetLength == 0 {
		return 0, fmt.Errorf("packet length is 0 bytes")
	}

	return int64(packetLength), nil
}
