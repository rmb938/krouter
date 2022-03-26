package v8

type Request struct {
	Topics                 []string
	AllowAutoTopicCreation bool
}

func (r *Request) Version() int16 {
	return Version
}
