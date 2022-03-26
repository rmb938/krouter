package v8

type Request struct {
	Topics []string
}

func (r *Request) Version() int16 {
	return Version
}
