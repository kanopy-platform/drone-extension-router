package defaults

import (
	"context"
	"encoding/json"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
)

type Config struct {
	Pipeline manifest.Pipeline `json:"pipeline,omitempty"`
}

type Defaults struct {
	config Config
}

func New(c Config) *Defaults {
	return &Defaults{config: c}
}

func (d *Defaults) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// decode pipeline resources
	resources, err := manifest.Decode(req.Config.Data)
	if err != nil {
		return nil, err
	}

	// merge defaults into user-defined resources
	for _, r := range resources {
		switch r.(type) {
		case *manifest.Pipeline:
			if err := merge(d.config.Pipeline, r); err != nil {
				return nil, err
			}
		}
	}

	// encode pipeline resources
	data, err := manifest.Encode(resources)
	if err != nil {
		return nil, err
	}

	return &drone.Config{
		Data: string(data),
	}, nil
}

func merge(defaults interface{}, user manifest.Resource) error {
	defaultBytes, err := json.Marshal(defaults)
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	userBytes, err = jsonpatch.MergePatch(defaultBytes, userBytes)
	if err != nil {
		return err
	}

	return json.Unmarshal(userBytes, user)
}
