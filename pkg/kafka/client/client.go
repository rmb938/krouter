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
	Broker *logical_broker.LogicalBroker

	log  logr.Logger
	conn net.Conn

	shutdown     bool
	shutdownChan chan struct{}

	requestChan  chan *netCodec.Packet
	responseChan chan struct {
		*netCodec.Packet
		int32
	}
}

func NewClient(log logr.Logger, broker *logical_broker.LogicalBroker, conn net.Conn) *Client {
	log = log.WithValues("from-address", conn.RemoteAddr().String())

	return &Client{
		Broker: broker,

		log:  log,
		conn: conn,

		shutdownChan: make(chan struct{}),

		requestChan: make(chan *netCodec.Packet, 1024),
		responseChan: make(chan struct {
			*netCodec.Packet
			int32
		}, 1024),
	}
}

func (c *Client) Run(packetProcessor *RequestPacketProcessor, requestHandler *RequestPacketHandler) {
	defer c.Close()

	// write out responses
	go func() {
		for c.shutdown == false {
			select {
			case <-c.shutdownChan:
				break
			case responseData := <-c.responseChan:
				correlationId := responseData.int32
				err := c.writePacket(responseData.Packet, correlationId)
				if err != nil {
					c.Shutdown()
					c.log.Error(err, "error writing packet")
					break
				}
			}
		}
	}()

	// handle requests
	go func() {
		for c.shutdown == false {
			select {
			case <-c.shutdownChan:
				break
			case requestPacket := <-c.requestChan:
				err := requestHandler.HandleRequest(c, requestPacket)
				if err != nil {
					c.Shutdown()
					c.log.Error(err, "error handling packet")
					break
				}
			}
		}
	}()

	// read in requests
	for c.shutdown == false {
		select {
		case <-c.shutdownChan:
			break
		default:
			packet, err := packetProcessor.ProcessPacket(c)
			if err != nil {
				c.Shutdown()
				c.log.Error(err, "error processing packet")
				break
			}
			c.requestChan <- packet
		}
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

	c.log.V(1).Info("Writing message", "correlation_id", correlationId)

	if c.shutdown == false {
		c.responseChan <- struct {
			*netCodec.Packet
			int32
		}{packet, correlationId}
	}

	return nil
}

func (c *Client) writePacket(packet *netCodec.Packet, correlationId int32) error {
	data, err := packet.Encode(correlationId)
	if err != nil {
		return fmt.Errorf("error encoding packet: %w", err)
	}

	c.log.V(1).Info("Writing packet", "key", packet.Key, "correlation_id", correlationId)
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

func (c *Client) Shutdown() {
	if c.shutdown {
		return
	}

	c.shutdown = true
	close(c.shutdownChan)
}

func (c *Client) Close() error {
	c.log.V(1).Info("Closed Connection")

	close(c.requestChan)
	close(c.responseChan)
	return c.conn.Close()
}
