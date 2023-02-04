package manifest

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

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
	ResourceData map[string]interface{} `yaml:",inline" json:",inline"`
}

func (r *resource) GetKind() Kind {
	return r.Kind
}

type Toleration struct {
	Key, Operator, Value, Effect string
}

type Pipeline struct {
	Kind         Kind                   `yaml:"kind" json:"kind"`
	Type         string                 `yaml:"type,omitempty" json:"type,omitempty"`
	Name         string                 `yaml:"name,omitempty" json:"name,omitempty"`
	NodeSelector map[string]string      `yaml:"node_selector,omitempty" json:"node_selector,omitempty"`
	Tolerations  []Toleration           `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
	ResourceData map[string]interface{} `yaml:",inline" json:",inline"`
}

func (p *Pipeline) GetKind() Kind {
	return KindPipeline
}

type Secret struct {
	Kind         Kind                   `yaml:"kind" json:"kind"`
	Name         string                 `yaml:"name,omitempty" json:"name,omitempty"`
	ResourceData map[string]interface{} `yaml:",inline" json:",inline"`
}

func (p *Secret) GetKind() Kind {
	return KindSecret
}

func Decode(data string) ([]Resource, error) {
	buf := bytes.NewBufferString(data)
	var resources []Resource
	dec := yaml.NewDecoder(buf)

	for {
		r := &resource{}

		err := dec.Decode(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch r.GetKind() {
		case KindPipeline:
			out, err := yaml.Marshal(r)
			if err != nil {
				return nil, err
			}

			p := &Pipeline{}
			if err := yaml.Unmarshal(out, p); err != nil {
				return nil, err
			}

			resources = append(resources, p)
		case KindSecret:
			out, err := yaml.Marshal(r)
			if err != nil {
				return nil, err
			}

			s := &Secret{}
			if err := yaml.Unmarshal(out, s); err != nil {
				return nil, err
			}

			resources = append(resources, s)
		default:
			resources = append(resources, r)
		}
	}

	return resources, nil
}

func Encode(resources []Resource) (string, error) {
	if len(resources) < 1 {
		return "", nil
	}

	buf := bytes.NewBuffer(nil)

	for idx, r := range resources {
		delim := "\n---\n"
		if idx == 0 {
			delim = "---\n"
		}

		resourceBytes, err := yaml.Marshal(r)
		if err != nil {
			return "", err
		}

		if _, err := buf.WriteString(delim); err != nil {
			return "", err
		}

		if _, err := buf.Write(resourceBytes); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
