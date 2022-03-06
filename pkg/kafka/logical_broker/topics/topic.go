package topics

type Topic struct {
	Name       string
	Partitions int32
	Config     map[string]string
	Enabled    bool
}

func (t *Topic) Clone() *Topic {
	newT := &Topic{
		Name:       t.Name,
		Partitions: t.Partitions,
		Config:     map[string]string{},
		Enabled:    t.Enabled,
	}

	for name, value := range t.Config {
		newT.Config[name] = value
	}

	return newT
}
