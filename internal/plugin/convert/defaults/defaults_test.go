package defaults_test

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/defaults"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	tests := []struct {
		desc    string
		config  defaults.Config
		request []manifest.Resource
		want    []manifest.Resource
	}{
		{
			desc: "test with empty defaults and request",
		},
		{
			desc:    "test without defaults",
			request: []manifest.Resource{&manifest.Pipeline{Kind: "pipeline"}},
			want:    []manifest.Resource{&manifest.Pipeline{Kind: "pipeline"}},
		},
		{
			desc: "test with defaults",
			config: defaults.Config{
				Pipeline: &manifest.Pipeline{NodeSelector: map[string]string{"d": "d", "test": "d"}},
			},
			request: []manifest.Resource{
				&manifest.Pipeline{Kind: "pipeline", NodeSelector: map[string]string{"r": "r", "test": "r"}},
			},
			want: []manifest.Resource{
				&manifest.Pipeline{Kind: "pipeline", NodeSelector: map[string]string{"r": "r", "d": "d", "test": "r"}},
			},
		},
		{
			desc: "test default node_selector and tolerations",
			config: defaults.Config{&manifest.Pipeline{
				NodeSelector: map[string]string{"instancegroup": "drone"},
				Tolerations: []manifest.Toleration{
					{Key: "dedicated", Operator: "Equal", Value: "drone", Effect: "NoSchedule"},
				},
			}},
			request: []manifest.Resource{
				&manifest.Pipeline{
					Kind:         "pipeline",
					NodeSelector: map[string]string{"instancegroup": "batch"},
					Tolerations: []manifest.Toleration{
						{Key: "dedicated", Operator: "Equal", Value: "batch", Effect: "NoSchedule"},
					},
				},
			},
			want: []manifest.Resource{
				&manifest.Pipeline{
					Kind:         "pipeline",
					NodeSelector: map[string]string{"instancegroup": "batch"},
					Tolerations: []manifest.Toleration{
						{Key: "dedicated", Operator: "Equal", Value: "batch", Effect: "NoSchedule"},
					},
				},
			},
		},
		{
			desc:   "test with multiple objects",
			config: defaults.Config{Pipeline: &manifest.Pipeline{Type: "test", Name: "test"}},
			request: []manifest.Resource{
				&manifest.Pipeline{Kind: "pipeline", Name: "user"},
				&manifest.Secret{Kind: "secret", Name: "user"},
			},
			want: []manifest.Resource{
				&manifest.Pipeline{Kind: "pipeline", Type: "test", Name: "user"},
				&manifest.Secret{Kind: "secret", Name: "user"},
			},
		},
	}

	for _, test := range tests {
		d := defaults.New(test.config)

		reqData, err := manifest.Encode(test.request)
		assert.NoError(t, err)

		req := &converter.Request{
			Config: drone.Config{
				Data: reqData,
			},
		}

		config, err := d.Convert(context.TODO(), req)
		assert.NoError(t, err)

		got, err := manifest.Decode(config.Data)
		assert.NoError(t, err)
		assert.Equal(t, test.want, got, test.desc)
	}
}
