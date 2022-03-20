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

	packetProcessor *client.RequestPacketProcessor
	requestHandler  *client.RequestPacketHandler
	listener        net.Listener

	broker *logical_broker.Broker
}

func (r *Router) ListenAndServe(listener, advertiseListener *net.TCPAddr, clusterID string, redisAddresses []string) error {
	r.Log.Info("Starting Router")

	r.packetProcessor = &client.RequestPacketProcessor{
		Log: r.Log.WithName("packet-processor"),
	}

	r.requestHandler = &client.RequestPacketHandler{
		Log: r.Log.WithName("request-handler"),
	}

	var err error
	r.listener, err = net.Listen("tcp", listener.String())
	if err != nil {
		return fmt.Errorf("error creating listener %w", err)
	}
	defer r.listener.Close()

	r.broker, err = logical_broker.InitBroker(r.Log, advertiseListener, clusterID, redisAddresses)
	if err != nil {
		return err
	}

	err = r.broker.InitClusters()
	if err != nil {
		return err
	}

	return r.serverLoop()
}

func (r *Router) Shutdown() error {
	// TODO: do something here to stop the serverLoop
	return nil
}

func (r *Router) serverLoop() error {
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			return err
		}

		log := r.Log.WithValues("from-address", conn.RemoteAddr().String())

		log.V(1).Info("Accepted Connection")

		c := client.NewClient(r.Log.WithName("client"), r.broker, conn)
		go c.Run(r.packetProcessor, r.requestHandler)
	}

	// TODO: wait for all clients to become idle (or timeout reached) then close brokers
	return r.broker.Close()
}
