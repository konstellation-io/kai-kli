package version_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/gookit/color.v1"

	"github.com/MakeNowJust/heredoc"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/kre/version"
	versioncmd "github.com/konstellation-io/kli/pkg/cmd/kre/version"

	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/konstellation-io/kli/mocks"
)

func TestVersionGetConfigCmd(t *testing.T) {
	s := newTestVersionSuite(t)
	versionConfig := &version.Config{
		Completed: true,
		Vars: []*version.ConfigVariable{
			{
				Key:   "key1",
				Value: "value1",
				Type:  version.ConfigVariableTypeVariable,
			},
			{
				Key:   "key2",
				Value: "value2",
				Type:  version.ConfigVariableTypeVariable,
			},
		},
	}
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().GetConfig("test-v1").Return(versionConfig, nil)

		return versioncmd.NewVersionCmd(f)
	})

	r.Run("version config test-v1").
		Containsf(heredoc.Doc(`
				TYPE     KEY
			1 VARIABLE key1
			2 VARIABLE key2

      [%s] Version config complete
		`), color.Success.Render("✔"))
}

func TestVersionGetConfigWithValuesCmd(t *testing.T) {
	s := newTestVersionSuite(t)
	versionConfig := &version.Config{
		Completed: true,
		Vars: []*version.ConfigVariable{
			{
				Key:   "key1",
				Value: "value1",
				Type:  version.ConfigVariableTypeVariable,
			},
			{
				Key:   "key2",
				Value: "value2",
				Type:  version.ConfigVariableTypeVariable,
			},
		},
	}
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().GetConfig("test-v1").Return(versionConfig, nil)

		return versioncmd.NewVersionCmd(f)
	})

	r.Run("version config test-v1 --show-values").
		Containsf(heredoc.Doc(`
				TYPE     KEY  VALUE
			1 VARIABLE key1 value1
			2 VARIABLE key2 value2

      [%s] Version config complete
		`), color.Success.Render("✔"))
}

func TestVersionSetConfigCmd(t *testing.T) {
	s := newTestVersionSuite(t)
	configVars := []version.ConfigVariableInput{
		{"key": "key1", "value": "value1"},
		{"key": "key2", "value": "value2"},
	}

	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().UpdateConfig("test-v1", configVars).Return(true, nil)

		return versioncmd.NewVersionCmd(f)
	})

	pair1 := fmt.Sprintf("%s=%s", configVars[0]["key"], configVars[0]["value"])
	pair2 := fmt.Sprintf("%s=%s", configVars[1]["key"], configVars[1]["value"])
	r.Runf("version config test-v1 --set %s --set %s", pair1, pair2).
		Containsf(heredoc.Doc(`
      [%s] Config completed for version 'test-v1'.
		`), color.Success.Render("✔"))
}

func TestVersionSetConfigErrorEdgeCasesCmd(t *testing.T) {
	_ = newTestVersionSuite(t)

	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		return versioncmd.NewVersionCmd(f)
	})

	type testCase struct {
		extraArgs     []string
		expectedError string
	}

	testCases := []testCase{
		{
			[]string{"--set", fmt.Sprintf("%s=%s", "key", "\"test")},
			`invalid argument "key=\"test" for "--set" flag: parse error on line 1, column 5: bare " in non-quoted-field`,
		},
		{
			[]string{"--set", fmt.Sprintf("%s=%s", "key", "\"test test\"")},
			`invalid argument "key=\"test test\"" for "--set" flag: parse error on line 1, column 5: bare " in non-quoted-field`,
		},
		{
			[]string{"--set", fmt.Sprintf("%q=%q", "key", `test`)},
			`invalid argument "\"key\"=\"test\"" for "--set" flag: parse error on line 1, column 5: extraneous or missing " in quoted-field`,
		},
	}

	for _, pair := range testCases {
		err := fmt.Errorf("%s", pair.expectedError)
		r.RunArgsE("version config test-v1", err, pair.extraArgs...)
	}
}

func TestVersionSetConfigEdgeCasesCmd(t *testing.T) {
	s := newTestVersionSuite(t)

	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil).AnyTimes()
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version).AnyTimes()
		s.mocks.version.EXPECT().UpdateConfig("test-v1", gomock.Any()).Return(true, nil).AnyTimes()

		return versioncmd.NewVersionCmd(f)
	})

	testCases := [][]string{
		{"--set", fmt.Sprintf("%s=%s", "key", "test test")},
		{"--set", fmt.Sprintf("%s=%s", "key", "test\n test")},
		{"--set", fmt.Sprintf("%s=%s", "key", `test`)},
		{"--set", fmt.Sprintf("%s=%s", "key", "test=")},
		// This test should work, but it doesn't.
		// {"--set", fmt.Sprintf("%s=%s", "key", "test \"test")},
	}

	for _, extraArgs := range testCases {
		r.RunArgs("version config test-v1", extraArgs...).
			Containsf(heredoc.Doc(`
      [%s] Config completed for version 'test-v1'.
		`), color.Success.Render("✔"))
	}
}

func TestVersionSetConfigFromEnvCmd(t *testing.T) {
	s := newTestVersionSuite(t)
	configVars := []version.ConfigVariableInput{
		{
			"key":   "TEST_VAR",
			"value": "test-v1\n222",
		},
		{
			"key":   "TEST_VAR2",
			"value": "test\" test",
		},
	}
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().UpdateConfig("test-v1", configVars).Return(true, nil)

		return versioncmd.NewVersionCmd(f)
	})

	err := os.Setenv("TEST_VAR", heredoc.Doc(`test-v1
                                                      222`))
	require.NoError(t, err)

	err = os.Setenv("TEST_VAR2", "test\" test")
	require.NoError(t, err)

	r.Run("version config test-v1 --set-from-env TEST_VAR --set-from-env TEST_VAR2").
		Containsf(heredoc.Doc(`
      [%s] Config completed for version 'test-v1'.
		`), color.Success.Render("✔"))
}

func TestVersionSetConfigFromEmptyEnvCmd(t *testing.T) {
	s := newTestVersionSuite(t)
	configVars := []version.ConfigVariableInput{
		{
			"key":   "TEST_VAR",
			"value": "",
		},
	}
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().UpdateConfig("test-v1", configVars).Return(false, nil)

		return versioncmd.NewVersionCmd(f)
	})

	err := os.Unsetenv("TEST_VAR")
	require.NoError(t, err)

	r.Run("version config test-v1 --set-from-env TEST_VAR").
		Containsf(heredoc.Doc(`
      [%s] Config updated for version 'test-v1'.
		`), color.Success.Render("✔"))
}

func TestVersionSetConfigFromFileCmd(t *testing.T) {
	s := newTestVersionSuite(t)
	configVars := []version.ConfigVariableInput{
		{
			"key":   "TEST_VAR",
			"value": "test-v16",
		},
	}
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		setupVersionConfig(t, f)

		f.EXPECT().KreClient("test").Return(s.mocks.kreClient, nil)
		s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
		s.mocks.version.EXPECT().UpdateConfig("test-v1", configVars).Return(true, nil)

		return versioncmd.NewVersionCmd(f)
	})

	d := os.TempDir()
	tempEnvFile, err := ioutil.TempFile(d, "TestVersionSetConfigFromFileCmd")
	require.NoError(t, err)

	defer os.RemoveAll(tempEnvFile.Name())

	_, _ = tempEnvFile.WriteString(heredoc.Doc(`
		TEST_VAR=test-v16
	`))

	r.Runf("version config test-v1 --set-from-file %s", tempEnvFile.Name()).
		Containsf(heredoc.Doc(`
      [%s] Config completed for version 'test-v1'.
		`), color.Success.Render("✔"))
}
