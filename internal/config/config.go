package config

import (
	"github.com/drone/drone-go/plugin/converter"
	"github.com/drone/drone-go/plugin/validator"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/defaults"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin/convert/pathschanged"
	"github.com/kanopy-platform/drone-extension-router/pkg/manifest"
	opa "github.com/kanopy-platform/drone-validation/plugin"
)

type (
	Config struct {
		Convert  Convert  `yaml:"convert"`
		Validate Validate `yaml:"validate"`
	}

	// convert
	Convert struct {
		Defaults     Defaults     `yaml:"defaults"`
		Pathschanged Pathschanged `yaml:"pathschanged"`
	}

	Defaults struct {
		Enable   bool               `yaml:"enable"`
		Pipeline *manifest.Pipeline `yaml:"pipeline,omitempty"`
	}

	Pathschanged struct {
		Enable bool `yaml:"enable"`
	}

	// validate
	Validate struct {
		OPA OPA `yaml:"opa"`
	}

	OPA struct {
		Enable     bool   `yaml:"enable"`
		PolicyPath string `yaml:"policyPath"`
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

func (c *Config) EnabledValidatePlugins() []validator.Plugin {
	plugins := []validator.Plugin{}

	if c.Validate.OPA.Enable {
		plugins = append(plugins, opa.New(opa.WithPolicyPath(c.Validate.OPA.PolicyPath)))
	}

	return plugins
}
