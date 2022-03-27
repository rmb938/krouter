package v5

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type GroupMember struct {
	MemberID string
	Metadata []byte
}

type Response struct {
	ErrCode      errors.KafkaError
	GenerationID int32
	ProtocolName string
	Leader       string
	MemberID     string
	Members      []GroupMember
}

func (r *Response) Version() int16 {
	return Version
}
