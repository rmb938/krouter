package v0

import "github.com/rmb938/krouter/pkg/kafka/message/impl/errors"

type DescribeGroupMembersResponse struct {
	MemberID         string
	ClientID         string
	ClientHost       string
	MemberMetadata   []byte
	MemberAssignment []byte
}

type DescribeGroupGroupResponse struct {
	ErrCode      errors.KafkaError
	GroupID      string
	GroupState   string
	ProtocolType string
	ProtocolData string
	Members      []DescribeGroupMembersResponse
}

type Response struct {
	Groups []DescribeGroupGroupResponse
}

func (r *Response) Version() int16 {
	return Version
}
