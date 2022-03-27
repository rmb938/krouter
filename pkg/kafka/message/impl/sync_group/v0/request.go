package v3

type GroupAssignments struct {
	MemberID   string
	Assignment []byte
}

type Request struct {
	GroupID      string
	GenerationID int32
	MemberID     string
	Assignments  []GroupAssignments
}

func (r *Request) Version() int16 {
	return Version
}
