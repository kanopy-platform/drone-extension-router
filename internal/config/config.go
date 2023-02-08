package config

import (
	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/pathschanged"
)

type (
	Config struct {
		Convert Convert `json:"convert"`
	}

	Convert struct {
		Pathschanged Pathschanged `json:"pathschanged"`
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

	if c.Convert.Pathschanged.Enable {
		plugins = append(plugins, pathschanged.New())
	}

	return plugins
}
