package v1

import (
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/create_acls/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {
	msg := &v1.Request{}

	var err error
	var creationsLength int32
	if creationsLength, err = reader.ArrayLength(); err != nil {
		return nil, err
	}

	for i := int32(0); i < creationsLength; i++ {
		creation := v1.RequestCreation{}

		if creation.ResourceType, err = reader.Int8(); err != nil {
			return nil, err
		}

		if creation.ResourceName, err = reader.String(); err != nil {
			return nil, err
		}

		if creation.ResourcePatternType, err = reader.Int8(); err != nil {
			return nil, err
		}

		if creation.Principal, err = reader.String(); err != nil {
			return nil, err
		}

		if creation.Host, err = reader.String(); err != nil {
			return nil, err
		}

		if creation.Operation, err = reader.Int8(); err != nil {
			return nil, err
		}

		if creation.PermissionType, err = reader.Int8(); err != nil {
			return nil, err
		}

		msg.Creations = append(msg.Creations, creation)
	}

	return msg, nil
}
