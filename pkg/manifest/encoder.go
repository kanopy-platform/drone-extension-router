package manifest

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

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
	enc := yaml.NewEncoder(buf)
	defer enc.Close()

	for _, r := range resources {
		if err := enc.Encode(r); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
