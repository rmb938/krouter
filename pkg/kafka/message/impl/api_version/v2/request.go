package v2

type Request struct {
}

func (r *Request) Version() int16 {
	return Version
}
