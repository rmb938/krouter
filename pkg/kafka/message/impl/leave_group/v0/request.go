package v0

type Request struct {
	GroupID  string
	MemberID string
}

func (r *Request) Version() int16 {
	return Version
}
