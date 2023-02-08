package config

import (
	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/pathschanged"
)

type (
	Config struct {
		Convert Convert `yaml:"convert"`
	}

	Convert struct {
		Pathschanged Pathschanged `yaml:"pathschanged"`
	}

	Pathschanged struct {
		Enable bool `yaml:"enable"`
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
