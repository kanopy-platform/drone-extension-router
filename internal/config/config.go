package config

import (
	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/defaults"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/pathschanged"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
)

type (
	Config struct {
		Convert Convert `json:"convert"`
	}

	Convert struct {
		Defaults     Defaults     `json:"defaults"`
		Pathschanged Pathschanged `json:"pathschanged"`
	}

	Defaults struct {
		Enable   bool               `json:"enable"`
		Pipeline *manifest.Pipeline `json:"pipeline,omitempty"`
	}

	Pathschanged struct {
		Enable bool `json:"enable"`
	}
)

func New() *Config {
	return &Config{}
}

func (c *Config) EnabledConvertPlugins() []converter.Plugin {
	plugins := []converter.Plugin{}

	if c.Convert.Defaults.Enable {
		plugins = append(plugins, defaults.New(defaults.Config{Pipeline: c.Convert.Defaults.Pipeline}))
	}

	if c.Convert.Pathschanged.Enable {
		plugins = append(plugins, pathschanged.New())
	}

	return plugins
}
