package cli

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-convert/internal/plugin"
	"github.com/kanopy-platform/drone-convert/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootCommand struct{}

func NewRootCommand() *cobra.Command {
	root := &RootCommand{}

	cmd := &cobra.Command{
		Use:               "drone-convert",
		PersistentPreRunE: root.persistentPreRunE,
		RunE:              root.runE,
	}

	cmd.PersistentFlags().String("log-level", "info", "Configure log level")
	cmd.PersistentFlags().String("listen-address", ":8080", "Server listen address")
	cmd.PersistentFlags().String("secret", "", "Token used to authenticate http requests to the extension")
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
		return fmt.Errorf("--secret flag is required")
	}

	if viper.GetBool("pathschanged-enabled") {
		convertPlugins = append(convertPlugins, plugin.NewPathsChanged())
	}

	log.Printf("Starting server on %s\n", addr)

	srv := server.New(secret, server.WithPluginRouter(
		plugin.NewRouter(
			plugin.WithConvertPlugins(convertPlugins...),
		),
	))

	return http.ListenAndServe(addr, srv)
}
