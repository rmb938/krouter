package v8

import (
	"fmt"

	v8 "github.com/rmb938/krouter/pkg/kafka/message/impl/metadata/v8"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v8.Request{}

	msg.Topics = nil
	topicsLength, err := reader.NullableArrayLength()
	if err != nil {
		return nil, err
	}
	// TODO: we shouldn't need *topicsLength > 0 however when running
	//  kafka-consumer-groups.sh it errors with Timed out waiting for a node assignment. Call: listOffsets on broker X
	//  according to docs a length of = 0 means return no topics while nil means return all topics
	//  but this seems to be wrong?
	if topicsLength != nil && *topicsLength > 0 {
		msg.Topics = make([]string, 0)
		for i := int32(0); i < *topicsLength; i++ {
			topicName, err := reader.String()
			if err != nil {
				return nil, err
			}

			msg.Topics = append(msg.Topics, topicName)
		}
	}

	if msg.AllowAutoTopicCreation, err = reader.Bool(); err != nil {
		return nil, fmt.Errorf("error bool for auto topic creation: %w", err)
	}
	if msg.IncludeClusterAuthorizedOperations, err = reader.Bool(); err != nil {
		return nil, fmt.Errorf("error bool for cluster authorized operations: %w", err)
	}
	if msg.IncludeTopicAuthorizedOperations, err = reader.Bool(); err != nil {
		return nil, fmt.Errorf("error bool for topic authoized operations: %w", err)
	}

	return msg, nil
}
