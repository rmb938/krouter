package v2

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type Response struct {
	ErrCode errors.KafkaError
	NodeID  int32
	Host    string
	Port    int32
}

func (r *Response) Version() int16 {
	return Version
}
