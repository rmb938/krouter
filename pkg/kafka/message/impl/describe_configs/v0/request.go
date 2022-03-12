package v0

type DescribeConfigResourceRequest struct {
	ResourceType      int8
	ResourceName      string
	ConfigurationKeys []string
}

type Request struct {
	Resources []DescribeConfigResourceRequest
}

func (r *Request) Version() int16 {
	return Version
}
