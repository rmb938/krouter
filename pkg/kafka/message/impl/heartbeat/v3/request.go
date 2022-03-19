package v3

type Request struct {
	GroupID         string
	GenerationID    int32
	MemberID        string
	GroupInstanceId *string
}

func (r *Request) Version() int16 {
	return Version
}
