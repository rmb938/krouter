package v1

import (
	"time"

	v1 "github.com/rmb938/krouter/pkg/kafka/message/impl/init_producer_id/v1"
	"github.com/rmb938/krouter/pkg/net/codec"
	"github.com/rmb938/krouter/pkg/net/message"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *codec.PackerReader) (message.Message, error) {

	msg := &v1.Request{}

	var err error
	msg.TransactionalID, err = reader.NullableString()
	if err != nil {
		return nil, err
	}

	transactionTimeoutMS, err := reader.Int32()
	if err != nil {
		return nil, err
	}
	msg.TransactionTimeoutDuration = time.Duration(transactionTimeoutMS) * time.Millisecond

	return msg, nil
}
