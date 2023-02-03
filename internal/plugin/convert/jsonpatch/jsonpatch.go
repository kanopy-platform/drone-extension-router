package jsonpatch

import (
	"context"
	"encoding/json"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	jpatch "github.com/evanphx/json-patch/v5"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"sigs.k8s.io/yaml"
)

type JsonPatch struct {
	Pipeline jpatch.Patch `json:"pipeline,omitempty"`
	Secret   jpatch.Patch `json:"secret,omitempty"`
}

func New(patches string) (*JsonPatch, error) {
	j := &JsonPatch{}

	if err := yaml.Unmarshal([]byte(patches), j); err != nil {
		return nil, err
	}

	return j, nil
}

func (j *JsonPatch) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	m, err := manifest.Decode(req.Config.Data)
	if err != nil {
		return nil, err
	}

	for _, resource := range m.Resources {
		resourceBytes, err := json.Marshal(resource)
		if err != nil {
			return nil, err
		}

		switch resource.GetKind() {
		case "pipeline":
			resourceBytes, err = j.patch(resourceBytes, j.Pipeline)
			if err != nil {
				return nil, err
			}
		case "secret":
			resourceBytes, err = j.patch(resourceBytes, j.Secret)
			if err != nil {
				return nil, err
			}
		}

		if err := json.Unmarshal(resourceBytes, resource); err != nil {
			return nil, err
		}
	}

	data, err := manifest.Encode(m)
	if err != nil {
		return nil, err
	}

	return &drone.Config{Data: data}, nil
}

func (j *JsonPatch) patch(data []byte, patch jpatch.Patch) ([]byte, error) {
	var err error

	opts := jpatch.NewApplyOptions()
	opts.EnsurePathExistsOnAdd = true

	data, err = patch.ApplyWithOptions(data, opts)
	if err != nil {
		return nil, err
	}

	return data, nil
}
