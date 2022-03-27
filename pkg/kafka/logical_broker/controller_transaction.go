package logical_broker

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kmsg"
)

func (c *Controller) FindTransactionCoordinator(transaction string) (*kmsg.FindCoordinatorResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	request := kmsg.NewPtrFindCoordinatorRequest()
	request.CoordinatorType = 1
	request.CoordinatorKey = transaction

	response, err := c.franzKafkaClient.Request(ctx, request)
	if err != nil {
		return nil, err
	}

	coordinatorResponse := response.(*kmsg.FindCoordinatorResponse)

	return coordinatorResponse, nil
}

func (c *Controller) InitProducer(transactionalID *string, transactionTimeoutDuration time.Duration) (*kmsg.InitProducerIDResponse, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	request := kmsg.NewPtrInitProducerIDRequest()
	request.TransactionalID = transactionalID
	request.TransactionTimeoutMillis = int32(transactionTimeoutDuration.Milliseconds())

	response, err := c.franzKafkaClient.Request(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*kmsg.InitProducerIDResponse), nil
}
