package v2

type Request struct {
	Key     string
	KeyType int8
}

func (r *Request) Version() int16 {
	return Version
}
