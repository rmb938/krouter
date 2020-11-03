package logical_broker

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"github.com/go-logr/logr"

	"github.com/rmb938/krouter/pkg/storage"
	"github.com/rmb938/krouter/pkg/storage/types"
)

type Broker struct {
	ClusterID string

	log logr.Logger

	id int32

	database storage.Database
}

func InitBroker(log logr.Logger, database storage.Database, id int32, advertiseListener *net.TCPAddr) (*Broker, error) {
	log = log.WithName("broker").WithValues("brokerID", id)

	ctx := context.TODO()

	log.V(-1).Info("Getting Cluster ID")
	clusterID, err := database.GetClusterID(ctx)
	if err != nil && err != storage.ErrNotFound {
		return nil, fmt.Errorf("error getting cluster id: %w", err)
	}

	if clusterID == nil {
		log.V(-1).Info("Generating Cluster ID")
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return nil, fmt.Errorf("error generating random cluster id: %w", err)
		}

		log.V(-1).Info("Creating Cluster ID")
		clusterID = func(s string) *string { return &s }(base64.StdEncoding.EncodeToString(b))
		err = database.SetClusterID(ctx, *clusterID)
		if err != nil {
			return nil, fmt.Errorf("error setting cluster id: %w", err)
		}
	} else {
		log.V(-1).Info("Found Cluster ID", "clusterID", clusterID)
	}

	log.V(-1).Info("Getting Broker")
	dbBroker, err := database.GetBroker(ctx, id)
	if err != nil && err != storage.ErrNotFound {
		return nil, fmt.Errorf("error getting broker: %w", err)
	}

	if dbBroker == nil {
		dbBroker = &types.Broker{
			ID:   id,
			Host: advertiseListener.IP.String(),
			Port: int32(advertiseListener.Port),
			Rack: nil,
		}

		log.V(-1).Info("Creating Broker")
		err := database.CreateBroker(ctx, dbBroker)
		if err != nil {
			return nil, fmt.Errorf("error creating broker: %w", err)
		}
	} else {
		dbBroker.Host = advertiseListener.IP.String()
		dbBroker.Port = int32(advertiseListener.Port)
		dbBroker.Rack = nil

		log.V(-1).Info("Updating Broker")
		err := database.UpdateBroker(ctx, dbBroker)
		if err != nil {
			return nil, fmt.Errorf("error updating broker: %w", err)
		}
	}

	broker := &Broker{
		log:       log,
		id:        id,
		ClusterID: *clusterID,
		database:  database,
	}

	go broker.leaderElection()

	return broker, nil
}

func (b *Broker) leaderElection() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		currentLeader, err := b.database.ControllerGetLeader(context.TODO())
		if err != nil {
			b.log.Error(err, "error getting controller leader")
			continue
		}

		if currentLeader == nil {
			b.log.Info("there is no leader")
			break
		}

		// leader is us but it can't be us since we didn't try to become leader yet
		if *currentLeader == b.id {
			b.log.Error(nil, "database says we are controller leader but we didn't try yet")
			continue
		}

		b.log.Info("current leader is", "leader", b.id)
		break
	}

	for ; true; <-ticker.C {
		currentLeader, err := b.database.ControllerGetLeader(context.TODO())
		if err != nil {
			b.log.Error(err, "error getting controller leader")
			continue
		}

		if currentLeader != nil && *currentLeader == b.id {
			continue
		}

		err = b.database.ControllerTryBecomeLeader(context.TODO(), b.id)
		if err != nil {
			b.log.Error(err, "error trying to become leader")
			continue
		}

		b.log.V(-1).Info("Became controller leader")
	}
}

func (b *Broker) GetBrokers(ctx context.Context) ([]types.Broker, error) {
	return b.database.GetBrokers(ctx)
}

func (b *Broker) GetControllerID(ctx context.Context) (*int32, error) {
	return b.database.ControllerGetLeader(ctx)
}
