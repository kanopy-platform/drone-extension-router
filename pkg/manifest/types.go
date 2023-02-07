package manifest

type Kind string

const (
	KindPipeline Kind = "pipeline"
	KindSecret   Kind = "secret"
)

type (
	Resource interface {
		GetKind() Kind
	}

	resource struct {
		Kind         Kind                   `yaml:"kind"`
		ResourceData map[string]interface{} `yaml:",inline"`
	}

	Toleration struct {
		Key, Operator, Value, Effect string
	}

	Pipeline struct {
		Kind         Kind                   `yaml:"kind"`
		Type         string                 `yaml:"type,omitempty"`
		Name         string                 `yaml:"name,omitempty"`
		NodeSelector map[string]string      `yaml:"node_selector,omitempty"`
		Tolerations  []Toleration           `yaml:"tolerations,omitempty"`
		ResourceData map[string]interface{} `yaml:",inline"`
	}

	Secret struct {
		Kind         Kind                   `yaml:"kind"`
		Name         string                 `yaml:"name,omitempty"`
		ResourceData map[string]interface{} `yaml:",inline"`
	}
)

func (r *resource) GetKind() Kind {
	return r.Kind
}

func (p *Pipeline) GetKind() Kind {
	return KindPipeline
}

func (p *Secret) GetKind() Kind {
	return KindSecret
}
