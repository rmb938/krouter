package logical_broker

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

func (c *Controller) APISetTopicPointer(topicName string, cluster string) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	record := kgo.KeySliceRecord([]byte(topicName), []byte(cluster))
	record.Topic = InternalTopicTopicPointers
	resp := c.franzKafkaClient.ProduceSync(ctx, record)
	if resp.FirstErr() != nil {
		return resp.FirstErr()
	}

	return nil
}

func (c *Controller) APIGetTopicPointer(topicName string) (*string, error) {
	// TODO: wait to be synced

	if pointer, ok := c.topicPointers.Load(topicName); ok {
		return &pointer, nil
	}

	return nil, nil
}

func (c *Controller) APIDeleteTopicPointer(topicName string) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	record := kgo.KeySliceRecord([]byte(topicName), nil)
	record.Topic = InternalTopicTopicPointers
	resp := c.franzKafkaClient.ProduceSync(ctx, record)
	if resp.FirstErr() != nil {
		return resp.FirstErr()
	}

	return nil
}
