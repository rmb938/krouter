package v3

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type Response struct {
	ThrottleDuration time.Duration
	ErrCode          errors.KafkaError
}

func (r *Response) Version() int16 {
	return Version
}
