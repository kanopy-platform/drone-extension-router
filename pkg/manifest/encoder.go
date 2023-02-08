package manifest

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

const separator = "\n---\n"

// Decode parses all YAML documents in the input string, unmarshals
// to the appropriate type, and outputs a slice of Resource objects.
func Decode(data string) ([]Resource, error) {
	if data == "" {
		return nil, nil
	}

	dec := yaml.NewDecoder(bytes.NewBufferString(data))

	var resources []Resource
	for {
		r := &Object{}

		switch err := dec.Decode(r); {
		case err == io.EOF:
			return resources, nil
		case err != nil:
			return nil, err
		}

		switch r.GetKind() {
		case KindPipeline:
			p := &Pipeline{}
			if err := remarshal(r, p); err != nil {
				return nil, err
			}

			resources = append(resources, p)
		case KindSecret:
			s := &Secret{}
			if err := remarshal(r, s); err != nil {
				return nil, err
			}

			resources = append(resources, s)
		default:
			resources = append(resources, r)
		}
	}
}

func remarshal(src, dst Resource) error {
	srcBytes, err := yaml.Marshal(src)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(srcBytes, dst)
}

// Encode marshals all of the input Resource objects into
// a single multi-document YAML string.
func Encode(resources []Resource) (string, error) {
	buf := bytes.NewBuffer(nil)
	for idx, r := range resources {
		if idx != 0 {
			if _, err := buf.WriteString(separator); err != nil {
				return "", err
			}
		}

		resourceBytes, err := yaml.Marshal(r)
		if err != nil {
			return "", err
		}

		if _, err := buf.Write(bytes.TrimSpace(resourceBytes)); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
