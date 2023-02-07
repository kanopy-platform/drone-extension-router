package manifest

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

const separator = "\n---\n"

// Decode parses all YAML documents in the input string, unmarshals
// to the appropriate type, and outputs a slice of Resource objects.
func Decode(data string) ([]Resource, error) {
	if data == "" {
		return nil, nil
	}

	docs := bytes.Split([]byte(data), []byte(separator))

	var resources []Resource
	for _, doc := range docs {
		r := &Object{}
		if err := yaml.Unmarshal(doc, r); err != nil {
			return nil, err
		}

		switch r.GetKind() {
		case KindPipeline:
			p := &Pipeline{}
			if err := yaml.Unmarshal(doc, p); err != nil {
				return nil, err
			}

			resources = append(resources, p)
		case KindSecret:
			s := &Secret{}
			if err := yaml.Unmarshal(doc, s); err != nil {
				return nil, err
			}

			resources = append(resources, s)
		default:
			resources = append(resources, r)
		}
	}

	return resources, nil
}

// Encode marshals all of the input Resource objects into
// a single multi-document YAML string.
func Encode(resources []Resource) (string, error) {
	if len(resources) < 1 {
		return "", nil
	}

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
