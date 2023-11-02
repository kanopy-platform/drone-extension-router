package cli

import (
	"net/http"
	"os"

	"github.com/kanopy-platform/drone-extension-router/internal/config"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin"
	"github.com/kanopy-platform/drone-extension-router/internal/server"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type RootCommand struct {
	config        *config.Config
	ConfigFile    string `split_words:"true"`
	ListenAddress string `split_words:"true" default:":8080"`
	LogLevel      string `split_words:"true" default:"info"`
	Secret        string `required:"true"`
}

func NewRootCommand() *cobra.Command {
	root := &RootCommand{config: config.New()}

	cmd := &cobra.Command{
		Use:               "drone-extension-router",
		PersistentPreRunE: root.persistentPreRunE,
		RunE:              root.runE,
	}

	cmd.AddCommand(newVersionCommand())
	return cmd
}

func (c *RootCommand) persistentPreRunE(cmd *cobra.Command, args []string) error {
	if err := envconfig.Process("drone", c); err != nil {
		return err
	}

	// set log level
	logLevel, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}

	log.SetLevel(logLevel)

	// read config file if specified
	if c.ConfigFile != "" {
		data, err := os.ReadFile(c.ConfigFile)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(data, c.config); err != nil {
			return err
		}
	}

	return nil
}

func (c *RootCommand) runE(cmd *cobra.Command, args []string) error {
	log.Printf("Starting server on %s\n", c.ListenAddress)

	pluginRouter := plugin.NewRouter(
		c.Secret,
		plugin.WithConvertPlugins(c.config.EnabledConvertPlugins()...),
		plugin.WithValidatePlugins(c.config.EnabledValidatePlugins()...),
		plugin.WithLogger(log.StandardLogger()),
	)

	return http.ListenAndServe(c.ListenAddress, server.New(pluginRouter))
}
