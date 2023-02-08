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
			desc:    "test with defaults",
			config:  "pipeline:\n  kind: pipeline\n  node_selector:\n    d: d\n    test: d\n  map:\n    d: d\n    test: d\n  list:\n  - test: d",
			request: "kind: pipeline\nnode_selector:\n  r: r\n  test: r\nmap:\n  r: r\n  test: r\nlist:\n- test: r",
			want:    "kind: pipeline\nnode_selector:\n    d: d\n    r: r\n    test: r\nlist:\n    - test: r\nmap:\n    d: d\n    r: r\n    test: r",
		},
		{
			desc:    "test default node_selector and tolerations",
			config:  "pipeline:\n  kind: pipeline\n  node_selector:\n    instancegroup: drone\n  tolerations:\n  - key: dedicated\n    operator: Equal\n    value: drone\n    effect: NoSchedule",
			request: "kind: pipeline\nnode_selector:\n  instancegroup: batch\ntolerations:\n- key: dedicated\n  operator: Equal\n  value: batch\n  effect: NoSchedule",
			want:    "kind: pipeline\nnode_selector:\n    instancegroup: batch\ntolerations:\n    - key: dedicated\n      operator: Equal\n      value: batch\n      effect: NoSchedule",
		},
		{
			desc:    "test with multiple objects",
			config:  "pipeline:\n  type: test\n  name: test",
			request: "kind: pipeline\nname: user\n---\nkind: secret\nname: user",
			want:    "kind: pipeline\ntype: test\nname: user\n---\nkind: secret\nname: user",
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
