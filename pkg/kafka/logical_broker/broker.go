package logical_broker

import (
	"net"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type Broker struct {
	AdvertiseListener *net.TCPAddr
	ClusterID         string

	log logr.Logger

	ephemeralID string
}

func InitBroker(log logr.Logger, advertiseListener *net.TCPAddr, clusterID string) (*Broker, error) {
	log = log.WithName("broker")

	broker := &Broker{
		AdvertiseListener: advertiseListener,
		ClusterID:         clusterID,

		log:         log,
		ephemeralID: uuid.Must(uuid.NewRandom()).String(),
	}

	return broker, nil
}
