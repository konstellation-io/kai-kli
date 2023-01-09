package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/konstellation-io/kli/cmd/factory"
	"github.com/konstellation-io/kli/pkg/cmd/args"
	"github.com/konstellation-io/kli/pkg/errors"
)

// NewActionCmd manages all the actions of the version command.
func NewActionCmd(f factory.CmdFactory, action string) *cobra.Command {
	log := f.Logger()
	cmd := &cobra.Command{
		Use:   action,
		Args:  args.ComposeArgsCheck(args.CheckServerFlag, cobra.ExactArgs(1)),
		Short: fmt.Sprintf("%s a version", cases.Title(language.English).String(action)),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, _ := cmd.Flags().GetString("server")
			c, err := f.KreClient(s)
			if err != nil {
				return err
			}

			comment, err := cmd.Flags().GetString("message")
			if err != nil {
				return err
			}

			versionName := args[0]
			version := c.Version()
			actionResult := ""
			switch action {
			case "start":
				err = version.Start(versionName, comment)
				actionResult = "Starting"
			case "stop":
				err = version.Stop(versionName, comment)
				actionResult = "Stopping"
			case "publish":
				err = version.Publish(versionName, comment)
				actionResult = "Publishing"
			case "unpublish":
				err = version.Unpublish(versionName, comment)
				actionResult = "Unpublishing"
			default:
				return errors.ErrUnknownVersionAction
			}
			if err != nil {
				return err
			}

			log.Success(fmt.Sprintf("%s version '%s'.", actionResult, versionName))
			return nil
		},
	}

	cmd.Flags().StringP("message", "m", "", fmt.Sprintf("Adds audit message to %s", action))
	_ = cmd.MarkFlagRequired("message")

	return cmd
}
