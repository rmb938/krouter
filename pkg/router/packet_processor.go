package router

import (
	"encoding/binary"
	"fmt"

	"github.com/go-logr/logr"
	apiVersionv0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"

	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/codec"
	"github.com/rmb938/krouter/pkg/kafka/message/handler"
	netCodec "github.com/rmb938/krouter/pkg/net/codec"
)

type PacketProcessor struct {
	Log logr.Logger
}

func (pp *PacketProcessor) processPacket(client *client.Client) error {
	log := pp.Log.WithValues("from-address", client.RemoteAddr().String())

	log.V(1).Info("Reading Packet Length")
	packetLength, err := pp.readPacketLength(client)
	if err != nil {
		return fmt.Errorf("error reading packet length: %w", err)
	}

	log.V(1).Info("Reading Packet Body", "length", packetLength)
	inPacketEncoded := make([]byte, packetLength)
	if _, err := client.ReadFull(inPacketEncoded); err != nil {
		return fmt.Errorf("error reading body of packet: %w", err)
	}

	inPacket, err := netCodec.DecodePacket(inPacketEncoded)
	if err != nil {
		return fmt.Errorf("error decoding packet: %w", err)
	}

	packetReader := netCodec.NewPacketReader(inPacket)

	log = log.WithValues("request_api_key", inPacket.ReqHeader.Key, "request_api_version", inPacket.ReqHeader.Version, "correlation_id", inPacket.ReqHeader.CorrelationId)

	decoderMap, ok := codec.MessageDecoderMapping[inPacket.ReqHeader.Key]
	if !ok {
		return fmt.Errorf("no decoder for request_api_key: %d", inPacket.ReqHeader.Key)
	}

	decoder, ok := decoderMap[inPacket.ReqHeader.Version]
	if !ok {
		log.Error(nil, "could not find a packet decoder, so returning unsupported version")
		err := client.WriteMessage(&apiVersionv0.Response{ErrCode: errors.UnsupportedVersion}, inPacket.ReqHeader.CorrelationId)
		if err != nil {
			return err
		}

		// don't error here so the kafka client can gracefully recover
		return nil
	}

	handlerMap, ok := handler.MessageHandlerMapping[inPacket.ReqHeader.Key]
	if !ok {
		return fmt.Errorf("no handler for packet: %d", inPacket.ReqHeader.Key)
	}

	inHandler, ok := handlerMap[inPacket.ReqHeader.Version]
	if !ok {
		return fmt.Errorf("no handler for packet version: %d", inPacket.ReqHeader.Version)
	}

	log.V(1).Info("decoding packet")
	message, err := decoder.Decode(packetReader)
	if err != nil {
		return fmt.Errorf("error decoding packet request_api_key: %d request_api_version: %d error: %w", inPacket.ReqHeader.Key, inPacket.ReqHeader.Version, err)
	}

	if err := packetReader.Close(); err != nil {
		return fmt.Errorf("error closing packet reader request_api_key: %d request_api_version: %d error: %w", inPacket.ReqHeader.Key, inPacket.ReqHeader.Version, err)
	}

	// TODO: figure out throttling stuff
	log.V(1).Info("handling packet")
	err = inHandler.Handle(client, pp.Log.WithValues("from-address", client.RemoteAddr().String()), message, inPacket.ReqHeader.CorrelationId)
	if err != nil {
		return fmt.Errorf("error handling packet: %w", err)
	}

	return nil
}

func (pp *PacketProcessor) readPacketLength(client *client.Client) (int64, error) {
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
