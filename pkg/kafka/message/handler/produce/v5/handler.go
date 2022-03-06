package v5

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/client"
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
	v5 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v5"
	"github.com/rmb938/krouter/pkg/kafka/records"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Handler struct {
}

func (h *Handler) Handle(client *client.Client, message message.Message, correlationId int32) error {
	request := message.(*v5.Request)

	response := &v5.Response{}

	for _, topicData := range request.TopicData {
		produceResponse := v5.ProduceResponse{
			Name: topicData.Name,
		}

		for _, partitionData := range topicData.PartitionData {

			// TODO: do something with this
			_, err := records.ParseRecordBatch(partitionData.Records)
			if err != nil {
				return err
			}

			partitionResponse := v5.PartitionResponse{}

			partitionResponse.Index = partitionData.Index
			partitionResponse.ErrCode = errors.None
			if request.TransactionalID != nil {
				// error if transaction because we don't support those yet
				partitionResponse.ErrCode = errors.TransactionIDAuthorizationFailed
			}

			// TODO: produce message to backend Kafka
			//  use the response from Kafka to fill these in as well

			partitionResponse.BaseOffset = 0             // TODO:  pull this from backend Kafka
			partitionResponse.LogAppendTime = time.Now() // TODO: pull this from backend Kafka
			partitionResponse.LogStartOffset = 0         // TODO:  pull this from backend Kafka

			produceResponse.PartitionResponses = append(produceResponse.PartitionResponses, partitionResponse)
		}

		response.Responses = append(response.Responses, produceResponse)
	}

	response.ThrottleDuration = 0

	return client.WriteMessage(response, correlationId)
}
