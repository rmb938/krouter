package topics

type Topic struct {
	Name       string
	Partitions int32
	Config     map[string]*string

	PrimaryCluster string
}

func (t *Topic) Clone() *Topic {
	newT := &Topic{
		Name:       t.Name,
		Partitions: t.Partitions,
		Config:     map[string]*string{},

		PrimaryCluster: t.PrimaryCluster,
	}

	for name, value := range t.Config {
		newT.Config[name] = value
	}

	return newT
}
