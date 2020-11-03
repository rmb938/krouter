package v9

type Brokers struct {
	ID   int32
	Host string
	Port int32
	Rack *string
}

type Partitions struct {
	Index           int32
	LeaderID        int32
	LeaderEpoch     int32
	ReplicaNodes    []int32
	ISRNodes        []int32
	OfflineReplicas []int32
}

type Topics struct {
	Name                      string
	Internal                  bool
	Partitions                []Partitions
	TopicAuthorizedOperations int32
}

type Response struct {
	Brokers                     []Brokers
	ClusterID                   string
	ControllerID                int32
	Topics                      []Topics
	ClusterAuthorizedOperations int32
}

func (r *Response) Version() int16 {
	return Version
}
