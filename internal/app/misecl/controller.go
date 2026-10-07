package misecl

import (
	"net"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/floffah/misecloud/api/atoc"
	"github.com/floffah/misecloud/internal/pkg/env"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ControllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Manage your mise cloud controller",
}

func InitController() {
	ControllerCmd.AddCommand(SetupCmd)
}

var SetupCmd = &cobra.Command{
	Use:   "setup <controller_dns>",
	Short: "Setup your mise cloud controller",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var confirm bool
		huh.NewConfirm().
			Title("Unconfigured controllers like this one are insecure by default. Continue?").
			Description("At the end of the process, this setup can optionally automatically configure gRPC TLS over Tailscale. In order to set this up you must enable enable HTTPS in your tailscale console. You can also do this manually later.").
			Affirmative("Continue").
			Negative("Cancel").
			Value(&confirm).
			WithButtonAlignment(lipgloss.Left).
			Run()

		argHost, argPort, err := net.SplitHostPort(args[0])
		if err != nil {
			argHost = args[0]
			argPort = env.CliViper.GetString(env.CliViperKeyPort)
		}

		address := net.JoinHostPort(argHost, argPort)
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		setupClient := atoc.NewControllerSetupClient(conn)

		isSetupResponse, err := setupClient.GetIsSetup(cmd.Context(), &emptypb.Empty{})
		if err != nil {
			panic(err)
		}

		if isSetupResponse.GetIsSetup() {
			var confirm bool
			huh.NewConfirm().
				Title("Controller is already set up. Continue?").
				Description("This may overwrite existing configuration and you could lose data. This action is irreversible. If you just want to configure the controller, please use `misecl controller` subcommands instead.").
				Affirmative("Continue").
				Negative("Cancel").
				Value(&confirm).
				WithButtonAlignment(lipgloss.Left).
				Run()

			if !confirm {
				//huh.NewNote().
				//	Title("Setup cancelled").
				//	Description("Controller setup has been cancelled. No changes have been made.").
				//	Run()
				return
			}
		}

		var setPort string = env.CliViper.GetString(env.CliViperKeyPort)
		huh.NewInput().
			Title("Controller Port").
			Description("Enter the port for the controller to listen on (default: 50051)").
			Value(&setPort).
			Run()

		_, err = setupClient.SetUp(cmd.Context(), &atoc.SetupRequest{Port: &setPort})
		if err != nil {
			panic(err)
		}

		huh.NewNote().
			Title("Setup complete").
			Description("Controller setup has been completed successfully.").
			Run()
	},
}
