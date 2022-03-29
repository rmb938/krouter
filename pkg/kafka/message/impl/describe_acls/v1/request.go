package v1

type Request struct {
	ResourceTypeFilter int8
	ResourceNameFilter *string
	PatternTypeFilter  int8
	PrincipalFilter    *string
	HostFilter         *string
	Operation          int8
	PermissionType     int8
}

func (r *Request) Version() int16 {
	return Version
}
