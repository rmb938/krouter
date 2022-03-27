package v2

type Request struct {
	Key string
}

func (r *Request) Version() int16 {
	return Version
}
