package config_test

import (
	"testing"

	"github.com/kanopy-platform/drone-extension-router/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	data := `---
convert:
  defaults:
    enable: true
    pipeline:
      node_selector:
        instancegroup: drone
      tolerations:
      - key: dedicated
        operator: Equal
        value: drone
        effect: NoSchedule
  pathschanged:
    enable: true
validate:
  opa:
    enable: true`

	c := config.New()
	assert.NoError(t, yaml.Unmarshal([]byte(data), c))

	// convert
	enabledConvert := c.EnabledConvertPlugins()
	assert.Len(t, enabledConvert, 2)

	assert.True(t, c.Convert.Defaults.Enable)
	assert.True(t, c.Convert.Pathschanged.Enable)
	assert.Equal(t, "drone", c.Convert.Defaults.Pipeline.NodeSelector["instancegroup"])
	assert.Equal(t, "drone", c.Convert.Defaults.Pipeline.Tolerations[0].Value)

	// validate
	enabledValidate := c.EnabledValidatePlugins()
	assert.Len(t, enabledValidate, 1)

	assert.True(t, c.Validate.OPA.Enable)
}
