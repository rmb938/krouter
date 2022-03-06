package logical_broker

import (
	"github.com/Shopify/sarama"
	"github.com/go-logr/logr"
)

type Controller struct {
	cluster *Cluster
}

func NewController(log logr.Logger, cluster *Cluster) (*Controller, error) {
	log = log.WithName("controller")

	return &Controller{cluster: cluster}, nil
}

func (c *Controller) findCoordinator(consumerGroup string) (*sarama.Broker, error) {
	return c.cluster.kafkaClient.Coordinator(consumerGroup)
}

func (c *Controller) JoinGroup(request *sarama.JoinGroupRequest) (*sarama.JoinGroupResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.GroupId)
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.JoinGroup(request)
}

func (c *Controller) SyncGroup(request *sarama.SyncGroupRequest) (*sarama.SyncGroupResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.GroupId)
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.SyncGroup(request)
}

func (c *Controller) LeaveGroup(request *sarama.LeaveGroupRequest) (*sarama.LeaveGroupResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.GroupId)
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.LeaveGroup(request)
}

func (c *Controller) OffsetFetch(request *sarama.OffsetFetchRequest) (*sarama.OffsetFetchResponse, error) {
	coordinatorBroker, err := c.findCoordinator(request.ConsumerGroup)
	if err != nil {
		return nil, err
	}

	return coordinatorBroker.FetchOffset(request)
}
