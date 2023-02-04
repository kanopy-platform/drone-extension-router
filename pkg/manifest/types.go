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
		Kind         Kind                   `yaml:"kind" json:"kind"`
		ResourceData map[string]interface{} `yaml:",inline" json:",inline"`
	}

	Toleration struct {
		Key, Operator, Value, Effect string
	}

	Pipeline struct {
		Kind         Kind                   `yaml:"kind" json:"kind"`
		Type         string                 `yaml:"type,omitempty" json:"type,omitempty"`
		Name         string                 `yaml:"name,omitempty" json:"name,omitempty"`
		NodeSelector map[string]string      `yaml:"node_selector,omitempty" json:"node_selector,omitempty"`
		Tolerations  []Toleration           `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
		ResourceData map[string]interface{} `yaml:",inline" json:",inline"`
	}

	Secret struct {
		Kind         Kind                   `yaml:"kind" json:"kind"`
		Name         string                 `yaml:"name,omitempty" json:"name,omitempty"`
		ResourceData map[string]interface{} `yaml:",inline" json:",inline"`
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
