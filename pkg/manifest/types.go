package manifest

type Kind string

const (
	KindPipeline Kind = "pipeline"
	KindSecret   Kind = "secret"
)

type Resource interface {
	GetKind() Kind
}

type resource struct {
	Kind         Kind                   `yaml:"kind" json:"kind"`
	Type         string                 `yaml:"type,omitempty" json:"type,omitempty"`
	Name         string                 `yaml:"name,omitempty" json:"name,omitempty"`
	ResourceData map[string]interface{} `yaml:",inline" json:",inline"`
}

func (r *resource) GetKind() Kind {
	return r.Kind
}

type (
	Pipeline struct {
		resource `yaml:",inline"`

		NodeSelector map[string]string `yaml:"node_selector,omitempty" json:"node_selector,omitempty"`
		Tolerations  []Toleration      `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
	}

	Toleration struct {
		Key, Operator, Value, Effect string
	}
)

func (p *Pipeline) GetKind() Kind {
	return KindPipeline
}

type Secret struct {
	resource `yaml:",inline"`
}

func (p *Secret) GetKind() Kind {
	return KindSecret
}
