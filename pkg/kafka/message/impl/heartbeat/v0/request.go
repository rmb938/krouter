package v0

type Request struct {
	GroupID      string
	GenerationID int32
	MemberID     string
}

func (r *Request) Version() int16 {
	return Version
}
