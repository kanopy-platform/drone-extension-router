package cli

import (
	"fmt"

	"github.com/kanopy-platform/drone-convert/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Build information for drone-convert",
		RunE: func(command *cobra.Command, args []string) error {
			fmt.Printf("%#v\n", version.Get())
			return nil
		},
	}
}
