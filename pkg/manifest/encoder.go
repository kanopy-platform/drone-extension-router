package manifest

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

func Decode(data string) ([]Resource, error) {
	dec := yaml.NewDecoder(bytes.NewBufferString(data))

	var resources []Resource
	for {
		r := &resource{}

		switch err := dec.Decode(r); {
		case err == io.EOF:
			return resources, nil
		case err != nil:
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
}

func Encode(resources []Resource) (string, error) {
	if len(resources) < 1 {
		return "", nil
	}

	buf := bytes.NewBuffer(nil)
	for idx, r := range resources {
		if idx != 0 {
			if _, err := buf.WriteString("\n---\n"); err != nil {
				return "", err
			}
		}

		resourceBytes, err := yaml.Marshal(r)
		if err != nil {
			return "", err
		}

		if _, err := buf.Write(resourceBytes); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
