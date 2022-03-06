package v2

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type Response struct {
	ThrottleDuration time.Duration
	ErrCode          errors.KafkaError
	ErrMessage       *string
	NodeID           int32
	Host             string
	Port             int32
}

func (r *Response) Version() int16 {
	return Version
}
