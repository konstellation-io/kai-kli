package krttools_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/cmd/krttools"
	"github.com/konstellation-io/kli/internal/krt"
	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/konstellation-io/kli/mocks"

	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var (
	sampleKrtString        = testhelpers.NewKrtBuilder().AsString()
	invalidSampleKrtString = testhelpers.NewKrtBuilder().WithWorkflows([]krt.Workflow{{
		Name:       "workflow",
		Entrypoint: "test",
		Nodes:      nil,
	}}).AsString()
)

type testKrtToolsSuite struct {
	ctrl  *gomock.Controller
	mocks suiteMocks
}

type suiteMocks struct {
	logger *mocks.MockLogger
}

func newTestKrtToolsSuite(t *testing.T) *testKrtToolsSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	return &testKrtToolsSuite{
		ctrl,
		suiteMocks{
			logger: logger,
		},
	}
}

func createKrtFile(t *testing.T, root, content, filename string) string {
	t.Helper()

	path := filepath.Join(root, filename)
	err := ioutil.WriteFile(path, []byte(content), 0600) //nolint:gocritic
	require.NoError(t, err)

	return path
}

func createTestKrtContent(t *testing.T, root string, files ...string) {
	t.Helper()

	for _, name := range files {
		name = filepath.Join(root, name)
		path := filepath.Dir(name)

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
			require.NoError(t, err)
		}

		f, err := os.Create(name)

		require.NoError(t, err)

		_ = f.Close()
	}
}

func TestNewKrtTools(t *testing.T) {
	s := newTestKrtToolsSuite(t)
	krtTools := krttools.NewKrtTools(s.mocks.logger).(*krttools.KrtTools)
	require.NotNil(t, krtTools)
}

func TestKrtTools_Validate(t *testing.T) {
	s := newTestKrtToolsSuite(t)

	tempDir, err := ioutil.TempDir("", "TestKrtTools_Validate")
	require.NoError(t, err)

	defer os.RemoveAll(tempDir)

	yamlFile := createKrtFile(t, tempDir, sampleKrtString, ".krt.yaml")
	krtTools := krttools.NewKrtTools(s.mocks.logger)
	err = krtTools.Validate(yamlFile)
	require.NoError(t, err)
}

func TestKrtTools_ValidateFailsWithKrtV1(t *testing.T) {
	s := newTestKrtToolsSuite(t)

	tempDir, err := ioutil.TempDir("", "TestKrtTools_Validate")
	require.NoError(t, err)

	defer os.RemoveAll(tempDir)

	yamlFile := createKrtFile(t, tempDir, invalidSampleKrtString, ".krt.yaml")
	krtTools := krttools.NewKrtTools(s.mocks.logger)
	err = krtTools.Validate(yamlFile)
	require.Error(t, err)
	require.ErrorContains(t, err, "error validating krt")
}

func TestKrtTools_Build(t *testing.T) {
	s := newTestKrtToolsSuite(t)
	files := []string{
		"src/py-test/main.py",
		"src/go-test/go.go",
		"docs/README.md",
		"metrics/dashboards/models.json",
		"metrics/dashboards/application.json",
	}

	// Create test dir structure
	tempDir, err := ioutil.TempDir("", "TestKrtTools_Build")
	require.NoError(t, err)

	defer os.RemoveAll(tempDir)

	createTestKrtContent(t, tempDir, files...)
	createKrtFile(t, tempDir, sampleKrtString, "krt.yaml")

	krtTools := krttools.NewKrtTools(s.mocks.logger)
	target := "test.krt"
	err = krtTools.Build(tempDir, target, "testversion")
	require.NoError(t, err)

	os.Remove(target)
}
