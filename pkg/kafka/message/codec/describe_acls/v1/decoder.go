package v1

import (
	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/describe_acls/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {
	msg := &v1.Request{}

	var err error
	if msg.ResourceTypeFilter, err = reader.Int8(); err != nil {
		return nil, err
	}

	if msg.ResourceNameFilter, err = reader.NullableString(); err != nil {
		return nil, err
	}

	if msg.PatternTypeFilter, err = reader.Int8(); err != nil {
		return nil, err
	}

	if msg.PrincipalFilter, err = reader.NullableString(); err != nil {
		return nil, err
	}

	if msg.HostFilter, err = reader.NullableString(); err != nil {
		return nil, err
	}

	if msg.Operation, err = reader.Int8(); err != nil {
		return nil, err
	}

	if msg.PermissionType, err = reader.Int8(); err != nil {
		return nil, err
	}

	return msg, nil
}
