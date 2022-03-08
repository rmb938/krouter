package v0

import "github.com/rmb938/krouter/pkg/kafka/message/impl/errors"

type Response struct {
	ErrCode errors.KafkaError
}

func (r *Response) Version() int16 {
	return Version
}
