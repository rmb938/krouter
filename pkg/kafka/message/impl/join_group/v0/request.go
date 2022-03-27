package v5

import "time"

type GroupProtocol struct {
	Name     string
	Metadata []byte
}

type Request struct {
	GroupID        string
	SessionTimeout time.Duration
	MemberID       string
	ProtocolType   string
	Protocols      []GroupProtocol
}

func (r *Request) Version() int16 {
	return Version
}
