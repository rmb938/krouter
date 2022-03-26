package client

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/message/codec"
	"github.com/rmb938/krouter/pkg/kafka/message/handler"
	apiVersionv0 "github.com/rmb938/krouter/pkg/kafka/message/impl/api_version/v0"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	netCodec "github.com/rmb938/krouter/pkg/net/codec"
)

type RequestPacketHandler struct {
	Log logr.Logger
}

func (rh *RequestPacketHandler) HandleRequest(client *Client, inPacket *netCodec.Packet) error {
	packetReader := netCodec.NewPacketReader(inPacket)

	log := rh.Log.WithValues("request_api_key", inPacket.ReqHeader.Key, "request_api_version", inPacket.ReqHeader.Version, "correlation_id", inPacket.ReqHeader.CorrelationId)

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
		return fmt.Errorf("no handler for packet %d version: %d", inPacket.ReqHeader.Key, inPacket.ReqHeader.Version)
	}

	log.V(1).Info("decoding packet")
	reqMessage, err := decoder.Decode(packetReader)
	if err != nil {
		return fmt.Errorf("error decoding packet request_api_key: %d request_api_version: %d error: %w", inPacket.ReqHeader.Key, inPacket.ReqHeader.Version, err)
	}

	if err := packetReader.Close(); err != nil {
		return fmt.Errorf("error closing packet reader request_api_key: %d request_api_version: %d error: %w", inPacket.ReqHeader.Key, inPacket.ReqHeader.Version, err)
	}

	// TODO: figure out throttling stuff
	log.Info("handling packet")
	respMessage, err := inHandler.Handle(client.Broker, log, reqMessage)
	if err != nil {
		return fmt.Errorf("error handling packet: %w", err)
	}

	if respMessage == nil {
		return nil
	}

	return client.WriteMessage(respMessage, inPacket.ReqHeader.CorrelationId)
}
