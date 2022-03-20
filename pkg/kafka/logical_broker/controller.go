package logical_broker

import (
	"github.com/go-logr/logr"
)

type Controller struct {
	cluster *Cluster
}

func NewController(log logr.Logger, cluster *Cluster) (*Controller, error) {
	log = log.WithName("controller")

	return &Controller{
		cluster: cluster,
	}, nil
}
