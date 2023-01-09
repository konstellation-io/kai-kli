package version_test

import (
	"fmt"
	"testing"

	"gopkg.in/gookit/color.v1"

	"github.com/MakeNowJust/heredoc"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/kre/version"
	"github.com/konstellation-io/kli/internal/config"
	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/konstellation-io/kli/mocks"
	cmd "github.com/konstellation-io/kli/pkg/cmd/kre/version"
)

type testVersionSuite struct {
	ctrl  *gomock.Controller
	mocks versionSuiteMocks
}

type versionSuiteMocks struct {
	kreClient *mocks.MockKreInterface
	version   *mocks.MockVersionInterface
}

func newTestVersionSuite(t *testing.T) *testVersionSuite {
	ctrl := gomock.NewController(t)

	return &testVersionSuite{
		ctrl,
		versionSuiteMocks{
			kreClient: mocks.NewMockKreInterface(ctrl),
			version:   mocks.NewMockVersionInterface(ctrl),
		},
	}
}

func setupVersionConfig(t *testing.T, f *mocks.MockCmdFactory) {
	cfg := f.Config()

	err := cfg.AddServer(config.ServerConfig{
		Name:     "test",
		URL:      "http://test.local",
		APIToken: "12346",
	})
	require.NoError(t, err)
	err = cfg.SetDefaultServer("test")
	require.NoError(t, err)
}

func TestVersionListCmd(t *testing.T) {
	s := newTestVersionSuite(t)
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().List().Return(version.List{
			{Name: "test-v1", Status: "STARTED"},
			{Name: "test-v2", Status: "STOPPED"},
		}, nil)

		return cmd.NewVersionCmd(f)
	})

	r.Run("version ls").
		Contains(heredoc.Doc(`
        NAME       STATUS
			1 test-v1 STARTED
			2 test-v2 STOPPED
		`))
}

func TestVersionStartNoMessageCmd(t *testing.T) {
	_ = newTestVersionSuite(t)
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		return cmd.NewVersionCmd(f)
	})

	r.RunE("version start test-v1", fmt.Errorf("required flag(s) \"message\" not set"))
}

func TestVersionStartCmd(t *testing.T) {
	comment := "test start comment"
	s := newTestVersionSuite(t)
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().Start("test-v1", comment).Return(nil)

		return cmd.NewVersionCmd(f)
	})

	r.Runf("version start test-v1 --message \"%s\"", comment).
		Containsf(heredoc.Doc(`
      [%s] Starting version 'test-v1'.
		`), color.Success.Render("✔"))
}

func TestVersionStopCmd(t *testing.T) {
	comment := "test stop comment"
	s := newTestVersionSuite(t)
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().Stop("test-v1", comment).Return(nil)

		return cmd.NewVersionCmd(f)
	})

	r.Runf("version stop test-v1 --message \"%s\"", comment).
		Containsf(heredoc.Doc(`
      [%s] Stopping version 'test-v1'.
		`), color.Success.Render("✔"))
}

func TestVersionPublishCmd(t *testing.T) {
	comment := "test publish comment"
	s := newTestVersionSuite(t)
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().Publish("test-v1", comment).Return(nil)

		return cmd.NewVersionCmd(f)
	})

	r.Runf("version publish test-v1 --message \"%s\"", comment).
		Containsf(heredoc.Doc(`
      [%s] Publishing version 'test-v1'.
		`), color.Success.Render("✔"))
}

func TestVersionUnpublishCmd(t *testing.T) {
	comment := "test unpublish comment"
	s := newTestVersionSuite(t)
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().Unpublish("test-v1", comment).Return(nil)

		return cmd.NewVersionCmd(f)
	})

	r.Runf("version unpublish test-v1 --message \"%s\"", comment).
		Containsf(heredoc.Doc(`
      [%s] Unpublishing version 'test-v1'.
		`), color.Success.Render("✔"))
}
