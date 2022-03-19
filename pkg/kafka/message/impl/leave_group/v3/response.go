package v3

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type ResponseMember struct {
	MemberID        string
	GroupInstanceID *string
	ErrCode         errors.KafkaError
}

type Response struct {
	ThrottleDuration time.Duration
	ErrCode          errors.KafkaError
	Members          []ResponseMember
}

func (r *Response) Version() int16 {
	return Version
}
