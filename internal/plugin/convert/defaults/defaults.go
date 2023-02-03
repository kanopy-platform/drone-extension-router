package defaults

import (
	"context"
	"encoding/json"

	"github.com/drone-runners/drone-runner-kube/engine/resource"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	dronemanifest "github.com/drone/runner-go/manifest"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

type Defaults struct {
	pipeline resource.Pipeline
}

type Option func(*Defaults)

func WithPipeline(p resource.Pipeline) Option {
	return func(d *Defaults) {
		d.pipeline = p
	}
}

func New(opts ...Option) *Defaults {
	d := &Defaults{}

	for _, o := range opts {
		o(d)
	}

	return d
}

func (d *Defaults) mergeDefaults(m *dronemanifest.Manifest) error {
	for idx, r := range m.Resources {
		userBytes, err := json.Marshal(r)
		if err != nil {
			return err
		}

		switch r.GetKind() {
		case dronemanifest.KindPipeline:
			defaultBytes, err := json.Marshal(d.pipeline)
			if err != nil {
				return err
			}

			userBytes, err = strategicpatch.StrategicMergePatch(defaultBytes, userBytes, d.pipeline)
			if err != nil {
				return err
			}
		default:
			return nil
		}

		if err := json.Unmarshal(userBytes, m.Resources[idx]); err != nil {
			return err
		}
	}

	return nil
}

func (d *Defaults) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	m, err := manifest.Decode(req.Config.Data)
	if err != nil {
		return nil, err
	}

	if err := d.mergeDefaults(m); err != nil {
		return nil, err
	}

	data, err := manifest.Encode(m)
	if err != nil {
		return nil, err
	}

	return &drone.Config{
		Data: string(data),
	}, nil
}
