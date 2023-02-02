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
		opts    []defaults.Option
		request string
		want    []dronemanifest.Resource
	}{
		{desc: "test with empty defaults and request"},
		{
			desc:    "test without defaults",
			request: "kind: pipeline\n",
			want: []dronemanifest.Resource{
				&resource.Pipeline{Kind: "pipeline"},
			},
		},
		{
			desc:    "test with defaults",
			opts:    []defaults.Option{defaults.WithPipeline(resource.Pipeline{NodeSelector: map[string]string{"d": "d", "test": "d"}})},
			request: "kind: pipeline\nnode_selector:\n  r: r\n  test: r",
			want: []dronemanifest.Resource{
				&resource.Pipeline{Kind: "pipeline", NodeSelector: map[string]string{"r": "r", "d": "d", "test": "r"}},
			},
		},
		{
			desc: "test default node_selector and tolerations",
			opts: []defaults.Option{
				defaults.WithPipeline(resource.Pipeline{
					NodeSelector: map[string]string{"instancegroup": "drone"},
					Tolerations: []resource.Toleration{
						{Key: "dedicated", Operator: "Equal", Value: "drone", Effect: "NoSchedule"},
					},
				}),
			},
			request: `---
kind: pipeline
node_selector:
  instancegroup: batch
tolerations:
- key: dedicated
  operator: Equal
  value: batch
  effect: NoSchedule
`,
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
			desc: "test with multiple objects",
			opts: []defaults.Option{defaults.WithPipeline(resource.Pipeline{Type: "test", Name: "test"})},
			request: `kind: pipeline
name: user
---
kind: secret
name: user
`,
			want: []dronemanifest.Resource{
				&resource.Pipeline{Kind: "pipeline", Type: "test", Name: "user"},
				&dronemanifest.Secret{Kind: "secret", Name: "user"},
			},
		},
	}

	for _, test := range tests {
		d := defaults.New(test.opts...)

		req := &converter.Request{
			Config: drone.Config{
				Data: test.request,
			},
		}

		config, err := d.Convert(context.TODO(), req)
		assert.NoError(t, err)

		got, err := manifest.Decode(config.Data)
		assert.Equal(t, test.want, got.Resources, test.desc)
	}
}
