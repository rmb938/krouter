package v1

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type ResponseResult struct {
	ErrCode    errors.KafkaError
	ErrMessage *string
}

type Response struct {
	ThrottleDuration time.Duration
	Results          []ResponseResult
}

func (r *Response) Version() int16 {
	return Version
}
