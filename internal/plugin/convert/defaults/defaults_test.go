package defaults_test

import (
	"context"
	"testing"

	"github.com/drone-runners/drone-runner-kube/engine/resource"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	dronemanifest "github.com/drone/runner-go/manifest"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/defaults"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	tests := []struct {
		desc    string
		config  defaults.Config
		request []dronemanifest.Resource
		want    []dronemanifest.Resource
	}{
		{
			desc: "test with empty defaults and request",
		},
		{
			desc:    "test without defaults",
			request: []dronemanifest.Resource{&resource.Pipeline{Kind: "pipeline"}},
			want:    []dronemanifest.Resource{&resource.Pipeline{Kind: "pipeline"}},
		},
		{
			desc: "test with defaults",
			config: defaults.Config{
				Pipeline: resource.Pipeline{NodeSelector: map[string]string{"d": "d", "test": "d"}},
			},
			request: []dronemanifest.Resource{
				&resource.Pipeline{Kind: "pipeline", NodeSelector: map[string]string{"r": "r", "test": "r"}},
			},
			want: []dronemanifest.Resource{
				&resource.Pipeline{Kind: "pipeline", NodeSelector: map[string]string{"r": "r", "d": "d", "test": "r"}},
			},
		},
		{
			desc: "test default node_selector and tolerations",
			config: defaults.Config{resource.Pipeline{
				NodeSelector: map[string]string{"instancegroup": "drone"},
				Tolerations: []resource.Toleration{
					{Key: "dedicated", Operator: "Equal", Value: "drone", Effect: "NoSchedule"},
				},
			}},
			request: []dronemanifest.Resource{
				&resource.Pipeline{
					Kind:         "pipeline",
					NodeSelector: map[string]string{"instancegroup": "batch"},
					Tolerations: []resource.Toleration{
						{Key: "dedicated", Operator: "Equal", Value: "batch", Effect: "NoSchedule"},
					},
				},
			},
			want: []dronemanifest.Resource{
				&resource.Pipeline{
					Kind:         "pipeline",
					NodeSelector: map[string]string{"instancegroup": "batch"},
					Tolerations: []resource.Toleration{
						{Key: "dedicated", Operator: "Equal", Value: "batch", Effect: "NoSchedule"},
					},
				},
			},
		},
		{
			desc:   "test with multiple objects",
			config: defaults.Config{Pipeline: resource.Pipeline{Type: "test", Name: "test"}},
			request: []dronemanifest.Resource{
				&resource.Pipeline{Kind: "pipeline", Name: "user"},
				&dronemanifest.Secret{Kind: "secret", Name: "user"},
			},
			want: []dronemanifest.Resource{
				&resource.Pipeline{Kind: "pipeline", Type: "test", Name: "user"},
				&dronemanifest.Secret{Kind: "secret", Name: "user"},
			},
		},
	}

	for _, test := range tests {
		d := defaults.New(test.config)

		reqData, err := manifest.Encode(&dronemanifest.Manifest{Resources: test.request})
		assert.NoError(t, err)

		req := &converter.Request{
			Config: drone.Config{
				Data: reqData,
			},
		}

		config, err := d.Convert(context.TODO(), req)
		assert.NoError(t, err)

		got, err := manifest.Decode(config.Data)
		assert.Equal(t, test.want, got.Resources, test.desc)
	}
}
