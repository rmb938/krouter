package router

import (
	"fmt"
	"net"

	"github.com/go-logr/logr"
	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/logical_broker"
)

type Router struct {
	log logr.Logger

	packetProcessor *client.RequestPacketProcessor
	requestHandler  *client.RequestPacketHandler
	listener        net.Listener

	Broker *logical_broker.Broker
}

func NewRouter(log logr.Logger, listener, advertiseListener *net.TCPAddr, clusterID string, redisAddresses []string) (*Router, error) {
	r := &Router{log: log}

	r.packetProcessor = &client.RequestPacketProcessor{
		Log: r.log.WithName("packet-processor"),
	}

	r.requestHandler = &client.RequestPacketHandler{
		Log: r.log.WithName("request-handler"),
	}

	var err error
	r.listener, err = net.Listen("tcp", listener.String())
	if err != nil {
		return nil, fmt.Errorf("error creating listener %w", err)
	}

	r.Broker, err = logical_broker.InitBroker(r.log, advertiseListener, clusterID, redisAddresses)
	if err != nil {
		return nil, err
	}

	err = r.Broker.InitClusters()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Router) ListenAndServe() error {
	r.log.Info("Starting Router")
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			return err
		}

		log := r.log.WithValues("from-address", conn.RemoteAddr().String())

		log.V(1).Info("Accepted Connection")

		c := client.NewClient(r.log.WithName("client"), r.Broker, conn)
		go c.Run(r.packetProcessor, r.requestHandler)
	}

	// TODO: wait for all clients to become idle (or timeout reached) then close brokers
	return r.Broker.Close()
}

func (r *Router) Shutdown() error {
	return r.listener.Close()
}
