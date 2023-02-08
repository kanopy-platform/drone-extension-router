package defaults_test

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/defaults"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDefaults(t *testing.T) {
	tests := []struct {
		desc    string
		config  string
		request string
		want    string
	}{
		{
			desc: "test with empty defaults and request",
		},
		{
			desc:    "test without defaults",
			request: "kind: pipeline",
			want:    "kind: pipeline",
		},
		{
			desc: "test with defaults",
			config: `
pipeline:
  node_selector:
    d: d
    test: d
  map:
    d: d
    test: d
  list:
  - test: d`,
			request: `
kind: pipeline
node_selector:
  r: r
  test: r
map:
  r: r
  test: r
list:
- test: r`,
			want: `kind: pipeline
node_selector:
    d: d
    r: r
    test: r
list:
    - test: r
map:
    d: d
    r: r
    test: r`,
		},
		{
			desc: "test default node_selector and tolerations",
			config: `
pipeline:
  kind: pipeline
  node_selector:
    instancegroup: drone
  tolerations:
  - key: dedicated
    operator: Equal
    value: drone
    effect: NoSchedule`,
			request: `
kind: pipeline
node_selector:
  instancegroup: batch
tolerations:
- key: dedicated
  operator: Equal
  value: batch
  effect: NoSchedule`,
			want: `kind: pipeline
node_selector:
    instancegroup: batch
tolerations:
    - key: dedicated
      operator: Equal
      value: batch
      effect: NoSchedule`,
		},
		{
			desc: "test with multiple objects",
			config: `
pipeline:
  type: test`,
			request: `
kind: pipeline
name: user
---
kind: secret
name: user`,
			want: `kind: pipeline
type: test
name: user
---
kind: secret
name: user`,
		},
	}

	for _, test := range tests {
		c := defaults.Config{}
		assert.NoError(t, yaml.Unmarshal([]byte(test.config), &c))

		req := &converter.Request{Config: drone.Config{Data: test.request}}

		config, err := defaults.New(c).Convert(context.TODO(), req)
		assert.NoError(t, err)

		assert.Equal(t, test.want, config.Data, test.desc)
	}
}
