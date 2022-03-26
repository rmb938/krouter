package v8

import (
	"github.com/rmb938/krouter/pkg/kafka/message/impl/errors"
)

type Brokers struct {
	ID   int32
	Host string
	Port int32
}

type Partitions struct {
	ErrCode      errors.KafkaError
	Index        int32
	LeaderID     int32
	ReplicaNodes []int32
	ISRNodes     []int32
}

type Topics struct {
	ErrCode    errors.KafkaError
	Name       string
	Partitions []Partitions
}

type Response struct {
	Brokers []Brokers
	Topics  []Topics
}

func (r *Response) Version() int16 {
	return Version
}
