package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/clientv3util"

	"github.com/rmb938/krouter/pkg/storage"
	"github.com/rmb938/krouter/pkg/storage/types"
)

const BrokersPath = "/brokers"

func (e *ETCD) GetBrokers(ctx context.Context) ([]types.Broker, error) {

	resp, err := e.client.Get(ctx, BrokersPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var brokers []types.Broker

	for _, kv := range resp.Kvs {
		broker := &types.Broker{}

		err = json.Unmarshal(kv.Value, broker)
		if err != nil {
			return nil, err
		}

		brokers = append(brokers, *broker)
	}

	return brokers, nil
}

func (e *ETCD) GetBroker(ctx context.Context, id int32) (*types.Broker, error) {

	brokerKey := fmt.Sprintf("%s/%d", BrokersPath, id)

	resp, err := e.client.KV.Txn(ctx).
		If(clientv3util.KeyExists(brokerKey)).
		Then(clientv3.OpGet(brokerKey)).
		Commit()

	if err != nil {
		return nil, err
	}

	if resp.Succeeded == false {
		return nil, storage.ErrNotFound
	}

	broker := &types.Broker{}

	kv := resp.Responses[0].GetResponseRange().Kvs[0]

	err = json.Unmarshal(kv.Value, broker)
	if err != nil {
		return nil, err
	}
	broker.StorageVersion = kv.Version

	return broker, nil
}

func (e *ETCD) CreateBroker(ctx context.Context, broker *types.Broker) error {

	data, err := json.Marshal(broker)
	if err != nil {
		return err
	}

	brokerKey := fmt.Sprintf("%s/%d", BrokersPath, broker.ID)

	resp, err := e.client.KV.Txn(ctx).
		If(clientv3util.KeyMissing(brokerKey)).
		Then(clientv3.OpPut(brokerKey, string(data))).
		Commit()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return storage.ErrExists
	}

	return nil
}

func (e *ETCD) UpdateBroker(ctx context.Context, broker *types.Broker) error {
	data, err := json.Marshal(broker)
	if err != nil {
		return err
	}

	brokerKey := fmt.Sprintf("%s/%d", BrokersPath, broker.ID)

	resp, err := e.client.KV.Txn(ctx).
		If(clientv3util.KeyExists(brokerKey), clientv3.Compare(clientv3.Version(brokerKey), "=", broker.StorageVersion)).
		Then(clientv3.OpPut(brokerKey, string(data))).
		Commit()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return storage.ErrNotFound
	}

	return nil
}

func (e *ETCD) DeleteBroker(ctx context.Context, id int32) error {
	panic("implement me")
}
