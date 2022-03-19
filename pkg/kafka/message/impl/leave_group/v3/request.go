package v3

type RequestMember struct {
	MemberID        string
	GroupInstanceID *string
}

type Request struct {
	GroupID string
	Members []RequestMember
}

func (r *Request) Version() int16 {
	return Version
}
