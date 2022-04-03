package models

import (
	"fmt"
	"time"
)

type BrokerEndpoint struct {
	Host string
	Port int32
}

type Broker struct {
	ID            int32
	Endpoint      BrokerEndpoint
	Rack          string
	LastHeartbeat time.Time
}

func (b *Broker) String() string {
	return fmt.Sprintf("%d-", b.ID)
}
