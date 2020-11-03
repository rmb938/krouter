package etcd

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/clientv3util"

	"github.com/rmb938/krouter/pkg/storage"
)

const ClusterPath = "/cluster"

func (e *ETCD) GetClusterID(ctx context.Context) (*string, error) {

	clusterIDKey := fmt.Sprintf("%s/id", ClusterPath)

	resp, err := e.client.KV.Txn(ctx).
		If(clientv3util.KeyExists(clusterIDKey)).
		Then(clientv3.OpGet(clusterIDKey)).
		Commit()

	if err != nil {
		return nil, err
	}

	if resp.Succeeded == false {
		return nil, storage.ErrNotFound
	}

	id := string(resp.Responses[0].GetResponseRange().Kvs[0].Value)

	return &id, nil
}

func (e *ETCD) SetClusterID(ctx context.Context, id string) error {

	clusterIDKey := fmt.Sprintf("%s/id", ClusterPath)

	resp, err := e.client.KV.Txn(ctx).
		If(clientv3util.KeyMissing(clusterIDKey)).
		Then(clientv3.OpPut(clusterIDKey, id)).
		Commit()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return storage.ErrExists
	}

	return nil
}
