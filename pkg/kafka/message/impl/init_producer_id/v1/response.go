package v1

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type Response struct {
	ThrottleDuration time.Duration
	ErrCode          errors.KafkaError
	ProducerID       int64
	ProducerEpoch    int16
}

func (r *Response) Version() int16 {
	return Version
}
