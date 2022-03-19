package v5

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type GroupMember struct {
	MemberID        string
	GroupInstanceID *string
	Metadata        []byte
}

type Response struct {
	ThrottleDuration time.Duration
	ErrCode          errors.KafkaError
	GenerationID     int32
	ProtocolName     string
	Leader           string
	MemberID         string
	Members          []GroupMember
}

func (r *Response) Version() int16 {
	return Version
}
