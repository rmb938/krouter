package etcd

import (
	"context"
	"fmt"
	"strconv"

	"go.etcd.io/etcd/clientv3/concurrency"
)

const ControllerPath = "/controller"

func (e *ETCD) ControllerTryBecomeLeader(ctx context.Context, brokerID int32) error {
	session, err := concurrency.NewSession(e.client)
	if err != nil {
		return err
	}

	controllerElectionPath := fmt.Sprintf("%s/election", ControllerPath)
	election := concurrency.NewElection(session, controllerElectionPath)

	// we want to try to become leader forever
	err = election.Campaign(context.Background(), strconv.Itoa(int(brokerID)))
	if err != nil && err != context.DeadlineExceeded {
		return err
	}

	return nil
}

func (e *ETCD) ControllerGetLeader(ctx context.Context) (*int32, error) {
	session, err := concurrency.NewSession(e.client)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	controllerElectionPath := fmt.Sprintf("%s/election", ControllerPath)
	election := concurrency.NewElection(session, controllerElectionPath)

	resp, err := election.Leader(ctx)
	if err != nil {
		if err == concurrency.ErrElectionNoLeader {
			return nil, nil
		}

		return nil, err
	}

	id64, err := strconv.ParseInt(string(resp.Kvs[0].Value), 10, 32)
	if err != nil {
		return nil, err
	}
	id := int32(id64)

	return &id, nil
}
