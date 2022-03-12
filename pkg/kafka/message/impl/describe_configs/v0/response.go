package v0

import (
	"time"

	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type DescribeConfigResultConfigResponse struct {
	Name      string
	Value     *string
	ReadOnly  bool
	Default   bool
	Sensitive bool
}

type DescribeConfigResultResponse struct {
	ErrCode      errors.KafkaError
	ErrMessage   *string
	ResourceType int8
	ResourceName string
	Configs      []DescribeConfigResultConfigResponse
}

type Response struct {
	ThrottleDuration time.Duration
	Results          []DescribeConfigResultResponse
}

func (r *Response) Version() int16 {
	return Version
}
