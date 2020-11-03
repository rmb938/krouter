package v9

type Request struct {
	Topics                             []string
	AllowAutoTopicCreation             bool
	IncludeClusterAuthorizedOperations bool
	IncludeTopicAuthorizedOperations   bool
}

func (r *Request) Version() int16 {
	return Version
}
