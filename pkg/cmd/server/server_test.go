package server_test

import (
	"testing"

	"gopkg.in/gookit/color.v1"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/internal/config"
	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/konstellation-io/kli/mocks"
	"github.com/konstellation-io/kli/pkg/cmd/server"
	"github.com/konstellation-io/kli/pkg/errors"
)

func TestServerListCmd(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		cfg := f.Config()

		err := cfg.AddServer(config.ServerConfig{
			Name:     "test",
			URL:      "http://test.local",
			APIToken: "12346",
		})
		require.NoError(t, err)

		return server.NewServerCmd(f)
	})

	r.Run("server ls").
		Contains(heredoc.Doc(`
			SERVER	URL
			test    http://test.local
		`))
}

func TestServerDefaultCmd(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		cfg := f.Config()

		err := cfg.AddServer(config.ServerConfig{
			Name:     "test",
			URL:      "http://test.local",
			APIToken: "12346",
		})
		require.NoError(t, err)

		return server.NewServerCmd(f)
	})

	r.Run("server ls").
		Contains(heredoc.Doc(`
			SERVER	URL
			test    http://test.local
		`))

	r.Run("server default test").
		Contains(heredoc.Doc(`
			SERVER	URL
			test*   http://test.local
		`)).
		Containsf(heredoc.Doc(`
			[%s] Server 'test' is now default.
		`), color.Success.Render("✔"))
}

func TestServerDefaultCmdError(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		return server.NewServerCmd(f)
	})

	r.RunE("server default unknown", errors.ErrUnknownServerName)
}

func TestServerAddCmd(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		return server.NewServerCmd(f)
	})
	r.Run("server add test http://test.local 12345").
		Containsf(heredoc.Doc(`
			[%s] Server 'test' added.
		`), color.Success.Render("✔")).
		Run("server ls").
		Contains(heredoc.Doc(`
			SERVER	URL
			test    http://test.local
		`))
}

func TestServerAddCmdErr(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		return server.NewServerCmd(f)
	})
	r.Run("server add test http://test.local 12345").
		Containsf(heredoc.Doc(`
			[%s] Server 'test' added.
		`), color.Success.Render("✔")).
		RunE("server add test http://test.local 12345", errors.ErrServerAlreadyExists)
}

func TestNoServerCmd(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		return server.NewServerCmd(f)
	})
	r.
		Run("server ls").
		Containsf(heredoc.Doc(`
			[%s] No servers found.
		`), color.Info.Render("ℹ"))
}
