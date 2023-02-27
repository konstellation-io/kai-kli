package config_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kli/internal/krt/config"
	"github.com/konstellation-io/kli/internal/krt/entity"
	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/konstellation-io/kli/mocks"
)

type testConfiguratorSuite struct {
	ctrl  *gomock.Controller
	mocks suiteMocks
}

type suiteMocks struct {
	logger   *mocks.MockLogger
	renderer *mocks.MockRenderer
}

func newTestBuilderSuite(t *testing.T) *testConfiguratorSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	return &testConfiguratorSuite{
		ctrl: ctrl,
		mocks: suiteMocks{
			logger:   logger,
			renderer: renderer,
		},
	}
}

func createFile(t *testing.T, filePath, content, filename string) {
	t.Helper()

	err := os.WriteFile(filepath.Join(filePath, filename), []byte(content), 0600)
	assert.NoError(t, err)
}

func TestConfigurator_Validate_ValidFile(t *testing.T) {
	s := newTestBuilderSuite(t)

	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	validKrtYaml := testhelpers.NewKrtBuilder().AsString()
	krtFileName := "krt.yaml"
	createFile(t, tmp, validKrtYaml, krtFileName)

	configurator := config.NewConfigurator(s.mocks.logger, s.mocks.renderer)

	err = configurator.Validate(filepath.Join(tmp, krtFileName))
	assert.NoError(t, err)
}

func TestConfigurator_Validate_NotValidFile(t *testing.T) {
	s := newTestBuilderSuite(t)

	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	validKrtYaml := testhelpers.NewKrtBuilder().WithVersion("_thisIs_not valid").AsString()
	krtFileName := "krt.yml"
	createFile(t, tmp, validKrtYaml, krtFileName)

	s.mocks.renderer.EXPECT().RenderValidationErrors(gomock.Len(1))

	configurator := config.NewConfigurator(s.mocks.logger, s.mocks.renderer)

	err = configurator.Validate(filepath.Join(tmp, krtFileName))
	assert.Error(t, err)
}

func TestConfigurator_Validate_KrtFileNoExists(t *testing.T) {
	s := newTestBuilderSuite(t)

	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	configurator := config.NewConfigurator(s.mocks.logger, s.mocks.renderer)

	err = configurator.Validate(filepath.Join(tmp, "krt.yaml"))
	assert.Error(t, err)
}

func TestConfigurator_UpdateVersion(t *testing.T) {
	s := newTestBuilderSuite(t)

	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	tmpInput, err := os.MkdirTemp(tmp, "input")
	assert.NoError(t, err)

	tmpOutput, err := os.MkdirTemp(tmp, "output")
	assert.NoError(t, err)

	validKrtYaml := testhelpers.NewKrtBuilder().AsString()
	krtFileName := "krt.yaml"
	createFile(t, tmpInput, validKrtYaml, krtFileName)

	configurator := config.NewConfigurator(s.mocks.logger, s.mocks.renderer)

	err = configurator.UpdateVersion("version-v1", tmpInput, tmpOutput)
	assert.NoError(t, err)

	outputKrt, err := entity.ParseFile(path.Join(tmpOutput, krtFileName))
	assert.NoError(t, err)

	assert.Equal(t, "version-v1", string(outputKrt.Version))
}

func TestConfigurator_UpdateVersion_InvalidVersion(t *testing.T) {
	s := newTestBuilderSuite(t)

	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	tmpInput, err := os.MkdirTemp(tmp, "input")
	assert.NoError(t, err)

	tmpOutput, err := os.MkdirTemp(tmp, "output")
	assert.NoError(t, err)

	validKrtYaml := testhelpers.NewKrtBuilder().AsString()
	krtFileName := "krt.yaml"
	createFile(t, tmpInput, validKrtYaml, krtFileName)

	configurator := config.NewConfigurator(s.mocks.logger, s.mocks.renderer)

	err = configurator.UpdateVersion("_Invalid version", tmpInput, tmpOutput)
	assert.ErrorContains(t, err, "invalid version name")
}

func TestConfigurator_UpdateVersion_InvalidInputPath(t *testing.T) {
	s := newTestBuilderSuite(t)

	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	tmpInput, err := os.MkdirTemp(tmp, "input")
	assert.NoError(t, err)

	tmpOutput, err := os.MkdirTemp(tmp, "output")
	assert.NoError(t, err)

	validKrtYaml := testhelpers.NewKrtBuilder().AsString()
	krtFileName := "krt.yaml"
	createFile(t, tmpInput, validKrtYaml, krtFileName)

	configurator := config.NewConfigurator(s.mocks.logger, s.mocks.renderer)

	err = configurator.UpdateVersion("version-v1", path.Join(tmpInput, krtFileName), tmpOutput)
	assert.ErrorContains(t, err, "input path should be a folder")
}

func TestConfigurator_UpdateVersion_InvalidInputPathWithNoYaml(t *testing.T) {
	s := newTestBuilderSuite(t)

	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	tmpInput, err := os.MkdirTemp(tmp, "input")
	assert.NoError(t, err)

	tmpOutput, err := os.MkdirTemp(tmp, "output")
	assert.NoError(t, err)

	configurator := config.NewConfigurator(s.mocks.logger, s.mocks.renderer)

	err = configurator.UpdateVersion("version-v1", tmpInput, tmpOutput)
	assert.ErrorContains(t, err, "krt.yaml not found in input folder")
}
