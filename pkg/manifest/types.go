package manifest

type Kind string

const (
	KindPipeline Kind = "pipeline"
	KindSecret   Kind = "secret"
)

type Resource interface {
	GetKind() Kind
}

// Object is used for intermediary marshalling of types that are
// not currently defined.
//
// Note that the `ResourceData` struct tag uses `,inline`, which
// ensures undefined fields are persisted when decoding/encoding.
type Object struct {
	ResourceData map[string]interface{} `yaml:",inline" json:",inline"`

	Kind Kind `yaml:"kind" json:"kind"`
}

func (r *Object) GetKind() Kind {
	return r.Kind
}

type (
	// Pipeline represents a Drone `pipeline` object.
	//
	// Note that the `ResourceData` struct tag uses `,inline`, which
	// ensures undefined fields are persisted when decoding/encoding.
	Pipeline struct {
		ResourceData map[string]interface{} `yaml:",inline" json:",inline"`

		Kind         Kind              `yaml:"kind" json:"kind"`
		Type         string            `yaml:"type,omitempty" json:"type,omitempty"`
		Name         string            `yaml:"name,omitempty" json:"name,omitempty"`
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

// Secret represents a Drone `secret` object.
//
// Note that the `ResourceData` struct tag uses `,inline`, which
// ensures undefined fields are persisted when decoding/encoding.
type Secret struct {
	ResourceData map[string]interface{} `yaml:",inline" json:",inline"`

	Kind Kind   `yaml:"kind" json:"kind"`
	Type string `yaml:"type,omitempty" json:"type,omitempty"`
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
}

func (p *Secret) GetKind() Kind {
	return KindSecret
}
