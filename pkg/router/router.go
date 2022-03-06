package router

import (
	"fmt"
	"net"

	"github.com/go-logr/logr"

	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
)

type Router struct {
	Log logr.Logger

	packetProcessor *PacketProcessor
	listener        net.Listener

	broker *logical_broker.Broker
}

func (r *Router) ListenAndServe(listener, advertiseListener *net.TCPAddr, clusterID string) error {
	r.Log.Info("Starting Router")

	r.packetProcessor = &PacketProcessor{
		Log: r.Log.WithName("packet-processor"),
	}

	var err error
	r.listener, err = net.Listen("tcp", listener.String())
	if err != nil {
		return fmt.Errorf("error creating listener %w", err)
	}

	defer func() {
		r.listener.Close()
	}()

	r.broker, err = logical_broker.InitBroker(r.Log, advertiseListener, clusterID)
	if err != nil {
		return err
	}

	return r.serverLoop()
}

func (r *Router) Shutdown() error {
	return r.listener.Close()
}

func (r *Router) serverLoop() error {
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			return err
		}

		log := r.Log.WithValues("from-address", conn.RemoteAddr().String())

		log.V(-1).Info("Accepted Connection")

		c := client.NewClient(r.Log.WithName("client"), r.broker, conn)
		go func() {
			defer c.Close()

			for {
				err := r.packetProcessor.processPacket(c)
				if err != nil {
					log.Error(err, "error processing packet")
					break
				}
			}
		}()
	}
}
