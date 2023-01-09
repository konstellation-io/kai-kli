package version_test

import (
	"testing"

	"gopkg.in/gookit/color.v1"

	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/konstellation-io/kli/mocks"
	versioncmd "github.com/konstellation-io/kli/pkg/cmd/kre/version"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func TestNewCreateCmd(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		s := newTestVersionSuite(t)
		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil).AnyTimes()
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version).AnyTimes()
		s.mocks.version.EXPECT().Create("test.krt").Return("1234", nil).AnyTimes()
		return versioncmd.NewVersionCmd(f)
	})

	r.Run("version create test.krt").
		Containsf(heredoc.Doc(`
      [%s] Upload KRT completed, version 1234 created.
		`), color.Success.Render("✔"))
}
