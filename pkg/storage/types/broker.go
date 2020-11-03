package types

type Broker struct {
	ID   int32   `json:"id"`
	Host string  `json:"host"`
	Port int32   `json:"port"`
	Rack *string `json:"rack,omitempty"`

	StorageVersion int64 `json:"-"`
}
