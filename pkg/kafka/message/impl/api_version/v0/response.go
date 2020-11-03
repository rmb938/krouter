package v0

type APIKey struct {
	Key        int16
	MinVersion int16
	MaxVersion int16
}

type Response struct {
	APIKeys []APIKey
}

func (r *Response) Version() int16 {
	return Version
}
