package storage

import (
	"golang.org/x/net/context"

	"github.com/rmb938/krouter/pkg/storage/types"
)

type Database interface {
	// Actions
	Close() error

	// Cluster Metadata
	GetClusterID(ctx context.Context) (*string, error)
	SetClusterID(ctx context.Context, id string) error

	// Controller
	ControllerTryBecomeLeader(ctx context.Context, brokerID int32) error
	ControllerGetLeader(ctx context.Context) (*int32, error)

	// Brokers
	GetBrokers(ctx context.Context) ([]types.Broker, error)
	GetBroker(ctx context.Context, id int32) (*types.Broker, error)
	CreateBroker(ctx context.Context, broker *types.Broker) error
	UpdateBroker(ctx context.Context, broker *types.Broker) error
	DeleteBroker(ctx context.Context, id int32) error
}
