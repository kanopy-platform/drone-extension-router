package manifest

import (
	"bytes"

	"github.com/drone-runners/drone-runner-kube/engine/resource"
	"github.com/drone/runner-go/manifest"
	dronemanifest "github.com/drone/runner-go/manifest"
	"sigs.k8s.io/yaml"
)

func init() {
	dronemanifest.Register(pipelineDriver)
}

func pipelineDriver(r *manifest.RawResource) (manifest.Resource, bool, error) {
	if r.Kind != resource.Kind {
		return nil, false, nil
	}

	out := new(resource.Pipeline)
	err := yaml.Unmarshal(r.Data, out)
	return out, true, err
}

func Decode(data string) (*dronemanifest.Manifest, error) {
	return dronemanifest.ParseString(data)
}

func Encode(m *dronemanifest.Manifest) (string, error) {
	if len(m.Resources) < 1 {
		return "", nil
	}

	s := bytes.NewBuffer(nil)

	for idx, r := range m.Resources {
		delim := "\n---\n"
		if idx == 0 {
			delim = "---\n"
		}

		b, err := yaml.Marshal(r)
		if err != nil {
			return "", err
		}

		if _, err := s.WriteString(delim); err != nil {
			return "", err
		}

		if _, err := s.Write(b); err != nil {
			return "", err
		}
	}

	return s.String(), nil
}
