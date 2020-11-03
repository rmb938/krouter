package client

import (
	"fmt"
	"io"
	"net"
	"reflect"

	"github.com/go-logr/logr"

	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
	"github.com/rmb938/krouter/pkg/kafka/message/codec"
	netCodec "github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Client struct {
	Broker *logical_broker.Broker

	log  logr.Logger
	conn net.Conn
}

func NewClient(log logr.Logger, broker *logical_broker.Broker, conn net.Conn) *Client {
	log = log.WithValues("from-address", conn.RemoteAddr().String())

	return &Client{
		Broker: broker,

		log:  log,
		conn: conn,
	}
}

func (c *Client) WriteMessage(msg message.Message, correlationId int32) error {

	messageType := reflect.TypeOf(msg).Elem()
	encoder, ok := codec.MessageEncoderMapping[messageType]
	if !ok {
		return fmt.Errorf("no encoder for message type %s: %s", messageType.PkgPath(), messageType.Name())
	}

	packet, err := encoder.Encode(msg)
	if err != nil {
		return fmt.Errorf("error encoding message: %w", err)
	}

	c.log.V(-1).Info("Writing message", "correlation_id", correlationId)
	return c.writePacket(packet, correlationId)
}

func (c *Client) writePacket(packet *netCodec.Packet, correlationId int32) error {
	data, err := packet.Encode(correlationId)
	if err != nil {
		return fmt.Errorf("error encoding packet: %w", err)
	}

	c.log.V(-1).Info("Writing packet", "key", packet.Key, "correlation_id", correlationId)
	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("error writting packet data: %w", err)
	}

	return nil
}

func (c *Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client) ReadFull(buf []byte) (n int, err error) {
	return io.ReadFull(c.conn, buf)
}

func (c *Client) Close() error {
	c.log.V(-1).Info("Closed Connection")
	return c.conn.Close()
}
