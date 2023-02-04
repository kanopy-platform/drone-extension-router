package defaults

import (
	"context"
	"encoding/json"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

type Config struct {
	Pipeline *manifest.Pipeline `json:"pipeline,omitempty"`
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
	for idx, r := range resources {
		userBytes, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}

		switch r.(type) {
		case *manifest.Pipeline:
			defaultBytes, err := json.Marshal(d.config.Pipeline)
			if err != nil {
				return nil, err
			}

			userBytes, err = strategicpatch.StrategicMergePatch(defaultBytes, userBytes, d.config.Pipeline)
			if err != nil {
				return nil, err
			}
		}

		if err := json.Unmarshal(userBytes, resources[idx]); err != nil {
			return nil, err
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
