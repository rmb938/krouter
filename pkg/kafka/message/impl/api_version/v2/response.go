package v2

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type APIKey struct {
	Key        int16
	MinVersion int16
	MaxVersion int16
}

type Response struct {
	ErrCode          errors.KafkaError
	APIKeys          []APIKey
	ThrottleDuration time.Duration
}

func (r *Response) Version() int16 {
	return Version
}
