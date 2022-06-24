package cli

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-extension-router/internal/plugin"
	"github.com/kanopy-platform/drone-extension-router/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootCommand struct{}

func NewRootCommand() *cobra.Command {
	root := &RootCommand{}

	cmd := &cobra.Command{
		Use:               "drone-extension-router",
		PersistentPreRunE: root.persistentPreRunE,
		RunE:              root.runE,
	}

	cmd.PersistentFlags().String("log-level", "info", "Configure log level")
	cmd.PersistentFlags().String("listen-address", ":8080", "Server listen address")
	cmd.PersistentFlags().Bool("pathschanged-enabled", false, "Enable pathschanged conversion extension")

	cmd.AddCommand(newVersionCommand())
	return cmd
}

func (c *RootCommand) persistentPreRunE(cmd *cobra.Command, args []string) error {
	// bind flags to viper
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("DRONE")
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// set log level
	logLevel, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return err
	}

	log.SetLevel(logLevel)

	return nil
}

func (c *RootCommand) runE(cmd *cobra.Command, args []string) error {
	convertPlugins := []converter.Plugin{}
	addr := viper.GetString("listen-address")

	secret := viper.GetString("secret")
	if secret == "" {
		return fmt.Errorf("DRONE_SECRET environment variable is required")
	}

	if viper.GetBool("pathschanged-enabled") {
		convertPlugins = append(convertPlugins, plugin.NewPathsChanged())
	}

	log.Printf("Starting server on %s\n", addr)

	pluginRouter := plugin.NewRouter(
		secret,
		plugin.WithConvertPlugins(convertPlugins...),
		plugin.WithLogger(log.StandardLogger()),
	)

	return http.ListenAndServe(addr, server.New(pluginRouter))
}
