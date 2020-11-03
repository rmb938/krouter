package etcd

import (
	"fmt"

	"go.etcd.io/etcd/clientv3"
)

type ETCD struct {
	client *clientv3.Client
}

func NewETCDStorage(endpoint string) (*ETCD, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{endpoint},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating etcd client: %w", err)
	}

	return &ETCD{client: etcdClient}, nil
}

func (e *ETCD) Close() error {
	return e.client.Close()
}
