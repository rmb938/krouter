package v1

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type ResponseACL struct {
	Principal      string
	Host           string
	Operation      int8
	PermissionType int8
}

type ResponseResource struct {
	ResourceType int8
	ResourceName string
	PatternType  int8
	ACLs         []ResponseACL
}

type Response struct {
	ThrottleDuration time.Duration
	ErrCode          errors.KafkaError
	ErrMessage       *string
	Resources        []ResponseResource
}

func (r *Response) Version() int16 {
	return Version
}
