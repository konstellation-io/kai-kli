package builder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/internal/builder"
	"github.com/konstellation-io/kli/internal/krt"
	"github.com/konstellation-io/kli/mocks"

	"github.com/konstellation-io/kli/internal/testhelpers"

	"github.com/stretchr/testify/require"
)

var sampleKrtString = testhelpers.NewKrtBuilder().AsString() //nolint:gochecknoglobals

type testBuilderSuite struct {
	ctrl  *gomock.Controller
	mocks suiteMocks
}

type suiteMocks struct {
	logger *mocks.MockLogger
}

func newTestBuilderSuite(t *testing.T) *testBuilderSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	return &testBuilderSuite{
		ctrl,
		suiteMocks{
			logger: logger,
		},
	}
}

func createKrtFile(t *testing.T, root, content, filename string) {
	t.Helper()

	err := ioutil.WriteFile(filepath.Join(root, filename), []byte(content), 0600) //nolint:gocritic

	require.NoError(t, err)
}

func createTestKrtContent(t *testing.T, root string, files ...string) {
	t.Helper()

	for _, name := range files {
		name = path.Join(root, name)
		filePath := path.Dir(name)
		_, err := os.Stat(filePath)

		if os.IsNotExist(err) {
			err = os.MkdirAll(filePath, 0755)
			require.NoError(t, err)
		}

		f, err := os.Create(name)

		defer func() { //nolint:gocritic
			_ = f.Close()
		}()
		require.NoError(t, err)
	}
}

func TestBuilder_Build(t *testing.T) {
	s := newTestBuilderSuite(t)
	files := []string{
		"src/py-test/main.py",
		"src/py-test/main.pyc",
		"src/go-test/go.go",
		"docs/README.md",
		"metrics/dashboards/models.json",
		"metrics/dashboards/application.json",
	}

	// Create test dir structure
	tmp, err := ioutil.TempDir("", "TestBuilder_Build")

	require.NoError(t, err)

	defer os.RemoveAll(tmp)

	createTestKrtContent(t, tmp, files...)
	createKrtFile(t, tmp, "*.pyc", ".krtignore")
	createKrtFile(t, tmp, sampleKrtString, "krt.yaml")

	validator := krt.NewCustomValidator(s.mocks.logger)
	b := builder.New(validator)
	err = b.Build(tmp, "test.krt")
	require.NoError(t, err)

	defer os.RemoveAll("test.krt")

	_, err = os.Stat("test.krt")
	require.NoError(t, err)
}

func TestBuilder_UpdateVersion(t *testing.T) {
	krtFileName := "krt.yaml"
	s := newTestBuilderSuite(t)
	tmp, err := ioutil.TempDir("", "TestBuilder_UpdateVersion")
	require.NoError(t, err)
	createKrtFile(t, tmp, sampleKrtString, krtFileName)

	krtFile := testhelpers.NewKrtBuilder().Build()

	validator := krt.NewCustomValidator(s.mocks.logger)
	updateVersion := "test1234"
	b := builder.New(validator)
	err = b.UpdateVersion(krtFile, path.Join(tmp, krtFileName), updateVersion)
	require.NoError(t, err)
}

func TestBuilder_UpdateVersionErrors(t *testing.T) {
	s := newTestBuilderSuite(t)
	tmp, err := ioutil.TempDir("", "TestBuilder_UpdateVersionErrors")
	require.NoError(t, err)

	validator := krt.NewCustomValidator(s.mocks.logger)
	b := builder.New(validator)
	updateVersion := "test1234"
	invalidVersion := "12345TEST"

	defer os.RemoveAll(tmp)

	krtFile := testhelpers.NewKrtBuilder().Build()
	// Test wrong  name
	err = b.UpdateVersion(krtFile, tmp, invalidVersion)
	require.EqualError(t, err, "invalid version name")

	// Test error opening
	krtFileName := "krt.yml"
	createKrtFile(t, tmp, sampleKrtString, krtFileName)
	err = os.Chmod(filepath.Join(tmp, krtFileName), 0444)
	require.NoError(t, err)
	err = b.UpdateVersion(krtFile, path.Join(tmp, krtFileName), updateVersion)
	require.EqualError(t, err, fmt.Sprintf("error while opening yaml file: open %v/krt.yml: permission denied", tmp))

	// Grant permissions to remove created files
	err = os.Chmod(filepath.Join(tmp, krtFileName), 0755)
	require.NoError(t, err)
}
