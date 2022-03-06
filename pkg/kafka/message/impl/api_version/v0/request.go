package v0

type Request struct {
}

func (r *Request) Version() int16 {
	return Version
}
