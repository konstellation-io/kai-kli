package testhelpers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/konstellation-io/kli/internal/logger"

	"github.com/golang/mock/gomock"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/internal/config"
	"github.com/konstellation-io/kli/mocks"
	"github.com/konstellation-io/kli/text"
)

type cmd func(f *mocks.MockCmdFactory) *cobra.Command

// NewRunner returns an Runner instance to execute a command and test output.
func NewRunner(t *testing.T, cmd cmd) Runner {
	t.Helper()
	ctrl := gomock.NewController(t)

	dir := SetupConfigDir(t)
	defer CleanConfigDir(t, dir)

	cfg := config.NewConfigTest()
	runner := &cmdRunner{t, nil, "", cfg, bytes.NewBufferString("")}

	f := mocks.NewMockCmdFactory(ctrl)
	f.EXPECT().Config().Return(cfg).AnyTimes()

	log := logger.NewWithWriter(runner.buf)

	f.EXPECT().Logger().AnyTimes().DoAndReturn(func() logger.Interface {
		return log
	})

	runner.root = cmd(f)

	return runner
}

// Runner interface for test command runner.
type Runner interface {
	Run(string) Runner
	RunE(string, error) Runner
	Runf(string, ...interface{}) Runner
	RunArgs(string, ...string) Runner
	RunArgsE(string, error, ...string) Runner

	Equal(string) Runner
	Contains(string) Runner
	Containsf(string, ...interface{}) Runner
	Empty() Runner
}

type cmdRunner struct {
	t    *testing.T
	root *cobra.Command
	out  string
	cfg  *config.Config
	buf  *bytes.Buffer
}

func (c *cmdRunner) NewBuffer() *bytes.Buffer {
	return c.buf
}

func (c *cmdRunner) Equal(expected string) Runner {
	expected = text.LinesTrim(expected)
	require.Equal(c.t, expected, c.out)

	return c
}

func (c *cmdRunner) Contains(expected string) Runner {
	actualClean := text.LinesTrim(c.out)
	expectedClean := text.LinesTrim(expected)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(actualClean, expectedClean, true)
	formatDiff := dmp.DiffPrettyText(diffs)

	require.Contains(c.t, actualClean, expectedClean, formatDiff)

	return c
}

func (c *cmdRunner) Containsf(expected string, args ...interface{}) Runner {
	expected = fmt.Sprintf(expected, args...)
	actualClean := text.LinesTrim(c.out)
	expectedClean := text.LinesTrim(expected)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(actualClean, expectedClean, true)
	formatDiff := dmp.DiffPrettyText(diffs)

	require.Contains(c.t, actualClean, expectedClean, formatDiff)

	return c
}

func (c *cmdRunner) Empty() Runner {
	require.Empty(c.t, strings.ReplaceAll(c.out, "\n", ""))

	return c
}

func (c *cmdRunner) Run(cmd string) Runner {
	assert := require.New(c.t)
	b := c.NewBuffer()

	c.out = ""
	c.root.SetOut(b)

	args := splitArgs(c.t, cmd)
	args = args[1:]

	c.root.SetArgs(args)

	err := c.root.Execute()
	assert.NoError(err)

	out, err := ioutil.ReadAll(b) //nolint:gocritic
	assert.NoError(err)

	c.out = "\n" + string(out)

	return c
}

func (c *cmdRunner) Runf(format string, args ...interface{}) Runner {
	return c.Run(fmt.Sprintf(format, args...))
}

func (c *cmdRunner) RunArgs(cmd string, extraArgs ...string) Runner {
	assert := require.New(c.t)
	b := c.NewBuffer()

	c.out = ""
	c.root.SetOut(b)

	args := splitArgs(c.t, cmd)
	args = args[1:]
	args = append(args, extraArgs...)

	c.root.SetArgs(args)

	err := c.root.Execute()
	assert.NoError(err)

	out, err := ioutil.ReadAll(b) //nolint:gocritic
	assert.NoError(err)

	c.out = "\n" + string(out)

	return c
}

func (c *cmdRunner) RunE(cmd string, expectedErr error) Runner {
	assert := require.New(c.t)
	b := c.NewBuffer()

	c.out = ""
	c.root.SetOut(b)

	args := strings.Split(cmd, " ")
	args = args[1:]

	c.root.SetArgs(args)

	actualErr := c.root.Execute()
	assert.EqualError(actualErr, expectedErr.Error())

	out, err := ioutil.ReadAll(b) //nolint:gocritic
	assert.NoError(err)

	c.out = "\n" + string(out)

	return c
}

func (c *cmdRunner) RunArgsE(cmd string, expectedErr error, extraArgs ...string) Runner {
	assert := require.New(c.t)
	b := c.NewBuffer()

	c.out = ""
	c.root.SetOut(b)

	args := strings.Split(cmd, " ")
	args = args[1:]
	args = append(args, extraArgs...)

	c.root.SetArgs(args)

	actualErr := c.root.Execute()
	assert.EqualError(actualErr, expectedErr.Error())

	out, err := ioutil.ReadAll(b) //nolint:gocritic
	assert.NoError(err)

	c.out = "\n" + string(out)

	return c
}

// This only accepts simple and well formatted arguments
// this commands will fail:
// "version start test-v1 --comment \"test test\""
// "version config --set key=\"test\""
// Use RunArgs and RunArgsE instead.
func splitArgs(t *testing.T, s string) []string {
	t.Helper()
	// Split string
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' ' // space

	fields, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}

	return fields
}
