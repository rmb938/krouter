package topics

type Topic struct {
	Name       string             `json:"name"`
	Cluster    string             `json:"cluster"`
	Partitions int32              `json:"partitions"`
	Config     map[string]*string `json:"config"`
}

func (t *Topic) Clone() *Topic {
	newT := &Topic{
		Name:       t.Name,
		Cluster:    t.Cluster,
		Partitions: t.Partitions,
		Config:     map[string]*string{},
	}

	for name, value := range t.Config {
		newT.Config[name] = value
	}

	return newT
}
