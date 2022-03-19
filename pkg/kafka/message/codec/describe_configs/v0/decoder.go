package v0

import (
	v0 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_configs/v0"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v0.Request{}

	var resourcesLength int32
	var err error
	if resourcesLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}

	for i := int32(0); i < resourcesLength; i++ {
		resource := v0.DescribeConfigResourceRequest{}

		if resource.ResourceType, err = reader.Int8(); err != nil {
			return nil, err
		}

		if resource.ResourceName, err = reader.String(); err != nil {
			return nil, err
		}

		resource.ConfigurationKeys = nil
		var configLength *int32
		if configLength, err = reader.NullableArrayLength(); err != nil {
			return nil, err
		}

		if configLength != nil {
			for i := int32(0); i < *configLength; i++ {
				var configKey string
				if configKey, err = reader.String(); err != nil {
					return nil, err
				}

				resource.ConfigurationKeys = append(resource.ConfigurationKeys, configKey)
			}
		}

		msg.Resources = append(msg.Resources, resource)
	}

	return msg, nil
}
