package v7

import (
	"time"

	v7 "github.com/rmb938/krouter/pkg/kafka/message/impl/produce/v7"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v7.Request{}

	var err error
	msg.TransactionalID, err = reader.NullableString()
	if err != nil {
		return nil, err
	}

	msg.ACKs, err = reader.Int16()
	if err != nil {
		return nil, err
	}

	timeoutMS, err := reader.Int32()
	if err != nil {
		return nil, err
	}

	msg.TimeoutDuration = time.Duration(timeoutMS) * time.Millisecond

	topicDataLength, err := reader.ArrayLength()
	if err != nil {
		return nil, err
	}

	if topicDataLength > 0 {
		for i := int32(0); i < topicDataLength; i++ {
			topicData := v7.TopicData{}

			topicData.Name, err = reader.String()
			if err != nil {
				return nil, err
			}

			partitionDataLength, err := reader.ArrayLength()
			if err != nil {
				return nil, err
			}

			if partitionDataLength > 0 {
				for i := int32(0); i < partitionDataLength; i++ {
					partitionData := v7.PartitionData{}

					partitionData.Index, err = reader.Int32()
					if err != nil {
						return nil, err
					}

					partitionData.Records, err = reader.Records()
					if err != nil {
						return nil, err
					}

					topicData.PartitionData = append(topicData.PartitionData, partitionData)
				}
			}

			msg.TopicData = append(msg.TopicData, topicData)
		}
	}

	return msg, nil
}
