package jsonpatch_test

import (
	"context"
	"testing"

	"github.com/drone-runners/drone-runner-kube/engine/resource"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	dronemanifest "github.com/drone/runner-go/manifest"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/jsonpatch"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestJsonPatch(t *testing.T) {
	config := `---
patchMap:
  pipelinePatch:
  - op: add
    path: /node_selector/instancegroup
    value: drone
  - op: add
    path: /tolerations/-
    value:
      key: dedicated
      operator: Equal
      value: drone
      effect: NoSchedule
  secretPatch:
  - op: add
    path: /name
    value: test
pipeline:
- pipelinePatch
secret:
- secretPatch
`

	j, err := jsonpatch.New(config)
	assert.NoError(t, err)

	req := `---
kind: pipeline
node_selector:
  test: test
tolerations:
- key: test
  operator: Equal
  value: test
  effect: NoSchedule

---
kind: secret
`

	conf, err := j.Convert(context.TODO(), &converter.Request{Config: drone.Config{Data: req}})
	assert.NoError(t, err)

	m, err := manifest.Decode(conf.Data)
	assert.NoError(t, err)

	assert.Equal(
		t,
		&resource.Pipeline{
			Kind:         "pipeline",
			NodeSelector: map[string]string{"test": "test", "instancegroup": "drone"},
			Tolerations: []resource.Toleration{
				{Key: "test", Operator: "Equal", Value: "test", Effect: "NoSchedule"},
				{Key: "dedicated", Operator: "Equal", Value: "drone", Effect: "NoSchedule"},
			},
		},
		m.Resources[0],
	)

	assert.Equal(t, &dronemanifest.Secret{Kind: "secret", Name: "test"}, m.Resources[1])
}
