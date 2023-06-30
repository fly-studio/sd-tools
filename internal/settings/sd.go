package settings

type SDOptions struct {
	Path       string            `yaml:"path" validate:"required,dir"`
	Models     map[string]string `yaml:"models"`
	ControlNet map[string]string `yaml:"control_net"`
	Scripts    string            `yaml:"scripts" validate:"required"`
	Embeddings string            `yaml:"embeddings" validate:"required"`
}

func (s *SDOptions) ValidPath(path string) bool {
	for _, v := range s.Models {
		if v == path {
			return true
		}
	}
	for _, v := range s.ControlNet {
		if v == path {
			return true
		}
	}

	if s.Scripts == path {
		return true
	}

	if s.Embeddings == path {
		return true
	}

	return false
}
