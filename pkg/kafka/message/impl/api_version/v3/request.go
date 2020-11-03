package v3

type Request struct {
	ClientSoftwareName    string
	ClientSoftwareVersion string
}

func (r *Request) Version() int16 {
	return Version
}
