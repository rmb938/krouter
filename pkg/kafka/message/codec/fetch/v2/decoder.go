package v2

import (
	"time"

	v2 "github.com/rmb938/krouter/pkg/kafka/message/impl/fetch/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v2.Request{}

	var err error
	if _, err = reader.Int32(); err != nil {
		return nil, err
	}

	var maxWaitMS int32
	if maxWaitMS, err = reader.Int32(); err != nil {
		return nil, err
	}

	msg.MaxWait = time.Duration(maxWaitMS) * time.Millisecond

	if msg.MinBytes, err = reader.Int32(); err != nil {
		return nil, err
	}

	var topicsLength int32
	if topicsLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}
	for i := int32(0); i < topicsLength; i++ {
		requestTopic := v2.FetchTopicRequest{}

		if requestTopic.Name, err = reader.String(); err != nil {
			return nil, err
		}

		var partitionsLength int32
		if partitionsLength, err = reader.ArrayLength(); err != nil {
			return nil, err
		}

		for i := int32(0); i < partitionsLength; i++ {
			requestPartition := v2.FetchPartitionRequest{}

			if requestPartition.Partition, err = reader.Int32(); err != nil {
				return nil, err
			}

			if requestPartition.FetchOffset, err = reader.Int64(); err != nil {
				return nil, err
			}

			if requestPartition.PartitionMaxBytes, err = reader.Int32(); err != nil {
				return nil, err
			}

			requestTopic.Partitions = append(requestTopic.Partitions, requestPartition)
		}

		msg.Topics = append(msg.Topics, requestTopic)
	}

	return msg, nil
}
