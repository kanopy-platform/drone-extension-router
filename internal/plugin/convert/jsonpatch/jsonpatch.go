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
	options  *jpatch.ApplyOptions
	Pipeline jpatch.Patch `json:"pipeline,omitempty"`
	Secret   jpatch.Patch `json:"secret,omitempty"`
}

func New(patches string) (*JsonPatch, error) {
	j := &JsonPatch{options: jpatch.NewApplyOptions()}
	j.options.EnsurePathExistsOnAdd = true

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
			resourceBytes, err = j.Pipeline.ApplyWithOptions(resourceBytes, j.options)
			if err != nil {
				return nil, err
			}
		case "secret":
			resourceBytes, err = j.Secret.ApplyWithOptions(resourceBytes, j.options)
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
