package v1

type RequestCreation struct {
	ResourceType        int8
	ResourceName        string
	ResourcePatternType int8
	Principal           string
	Host                string
	Operation           int8
	PermissionType      int8
}

type Request struct {
	Creations []RequestCreation
}

func (r *Request) Version() int16 {
	return Version
}
