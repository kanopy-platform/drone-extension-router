package config

import (
	"github.com/drone-runners/drone-runner-kube/engine/resource"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/defaults"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/pathschanged"
)

type (
	Config struct {
		Defaults     Defaults     `json:"defaults"`
		Pathschanged Pathschanged `json:"pathschanged"`
	}

	Defaults struct {
		Enable   bool               `json:"enable"`
		Pipeline *resource.Pipeline `json:"pipeline,omitempty"`
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

	if c.Defaults.Enable {
		plugins = append(plugins, defaults.New(defaults.Config{Pipeline: c.Defaults.Pipeline}))
	}

	if c.Pathschanged.Enable {
		plugins = append(plugins, pathschanged.New())
	}

	return plugins
}
