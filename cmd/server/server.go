package server

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

const (
	_serverFlag = "server"
)

// NewServerCmd creates a new command to handle 'server' subcommands.
func NewServerCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <command> [opts...]",
		Short: "Manage servers for kli",
		Example: heredoc.Doc(`
			$ kli server ls
			$ kli server add <server_name> <server_url>
			$ kli server remove <server_name>
			$ kli server default <server_name>
			$ kli server login <server_name>
			$ kli server login-cli <server_name> --username <username> --password <password> ` +
			`--client-id <client-id> --client-secret <client-secret>
			$ kli server logout <server_name>
			$ kli server get-token my-server
		`),
	}

	cmd.AddCommand(
		NewListCmd(logger),
		NewDefaultCmd(logger),
		NewAddCmd(logger),
		NewRemoveCmd(logger),
		NewLoginCmd(logger),
		NewLoginCLICmd(logger),
		NewLogoutCmd(logger),
		NewGetTokenCmd(logger),
	)

	return cmd
}

func getServerOrDefault(cmd *cobra.Command, logger logging.Interface) (*configuration.Server, error) {
	serverName, err := cmd.Flags().GetString(_serverFlag)
	if err != nil {
		serverName = ""
	}

	configService := configuration.NewKaiConfigService(logger)

	kaiConfig, err := configService.GetConfiguration()
	if err != nil {
		return nil, err
	}

	srv, err := kaiConfig.GetServerOrDefault(serverName)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
