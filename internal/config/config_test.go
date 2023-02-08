package config_test

import (
	"testing"

	"github.com/kanopy-platform/drone-extension-router/internal/config"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestConfig(t *testing.T) {
	data := `---
convert:
  pathschanged:
    enable: true`

	c := config.New()
	assert.NoError(t, yaml.Unmarshal([]byte(data), c))

	enabled := c.EnabledConvertPlugins()
	assert.Len(t, enabled, 2)

	assert.True(t, c.Convert.Pathschanged.Enable)
}
