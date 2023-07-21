package server

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

var (
	ErrServerNotFound = errors.New("server not found")
)

const (
	_serverFlag = "server"
)

// NewLogoutCmd creates a new command to log out from an existing server.
func NewLogoutCmd(logger logging.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use: "logout [--server <server_name>]",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Short: "logout from an existing server",
		Example: `
    $ kli server logout my-server
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName, err := getServerName(cmd, logger)
			if err != nil {
				return ErrServerNotFound
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = server.NewServerHandler(logger, r).Logout(serverName)
			if err != nil {
				return err
			}

			logger.Success(fmt.Sprintf("Logged out from %q.", serverName))

			return nil
		},
	}

	cmd.Flags().StringP(_serverFlag, "s", "", "server name")

	return cmd
}

func getServerName(cmd *cobra.Command, logger logging.Interface) (string, error) {
	serverName, err := cmd.Flags().GetString(_serverFlag)
	if err == nil && serverName != "" {
		return serverName, nil
	}

	configService := configuration.NewKaiConfigService(logger)

	kaiConfig, err := configService.GetConfiguration()
	if err != nil {
		return "", err
	}

	defaultServer, err := kaiConfig.GetDefaultServer()
	if err != nil {
		return "", err
	}

	return defaultServer.Name, nil
}
