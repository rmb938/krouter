package v1

import "time"

type Request struct {
	TransactionalID            *string
	TransactionTimeoutDuration time.Duration
}

func (r *Request) Version() int16 {
	return Version
}
