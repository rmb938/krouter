package v0

type Request struct {
	Groups []string
}

func (r *Request) Version() int16 {
	return Version
}
