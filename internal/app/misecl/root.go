package misecl

import (
	"github.com/floffah/misecloud/internal/pkg/env"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "mise",
	Short:   "Command line interface for interacting with your misecloud",
	Version: "0.0.1",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := env.InitCliViper()
		if err != nil {
			return err
		}

		return nil
	},
}

func InitCli() {
	RootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")

	RootCmd.AddCommand(ControllerCmd)
	InitController()
}
