package kre_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kli/api/kre/runtime"
	"github.com/konstellation-io/kli/api/kre/version"
	"github.com/konstellation-io/kli/internal/kre"
	"github.com/konstellation-io/kli/mocks"
)

type testInteractorSuite struct {
	ctrl  *gomock.Controller
	mocks suiteMocks
}

type suiteMocks struct {
	logger    *mocks.MockLogger
	kreClient *mocks.MockKreInterface
	version   *mocks.MockVersionInterface
	runtime   *mocks.MockRuntimeInterface
	renderer  *mocks.MockRenderer
}

func newTestInteractorSuite(t *testing.T) *testInteractorSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	kreClient := mocks.NewMockKreInterface(ctrl)
	versionMock := mocks.NewMockVersionInterface(ctrl)
	runtimeMock := mocks.NewMockRuntimeInterface(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)

	return &testInteractorSuite{
		ctrl: ctrl,
		mocks: suiteMocks{
			logger:    logger,
			kreClient: kreClient,
			version:   versionMock,
			runtime:   runtimeMock,
			renderer:  renderer,
		},
	}
}

func TestInteractor_ListVersions(t *testing.T) {
	s := newTestInteractorSuite(t)

	versionRuntime := version.Version{
		Name:   "test-version",
		Status: "STARTED",
	}
	versionList := version.List([]version.Version{versionRuntime})

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-runtime").Return(versionList, nil)
	s.mocks.renderer.EXPECT().RenderVersions(versionList).Return()

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-runtime")
	assert.NoError(t, err)
}

func TestInteractor_ListVersions_NoVersionsInRuntime(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-runtime").Return(nil, nil)
	s.mocks.renderer.EXPECT().RenderVersions(nil).Return()

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-runtime")
	assert.NoError(t, err)
}

func TestInteractor_ListVersions_RequestFail(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-runtime").Return(nil, errors.New("client error"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-runtime")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error listing versions")
}

func TestInteractor_CreateVersion(t *testing.T) {
	s := newTestInteractorSuite(t)
	pathToKrt := "path/to/krt"
	expectedRuntime := "test-runtime"

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().Create(expectedRuntime, pathToKrt).Return("version-name", nil)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateVersion(expectedRuntime, pathToKrt)
	assert.NoError(t, err)
}

func TestInteractor_CreateVersion_ClientError(t *testing.T) {
	s := newTestInteractorSuite(t)
	pathToKrt := "path/to/krt"
	expectedRuntime := "test-runtime"

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().Create(expectedRuntime, pathToKrt).Return("", errors.New("error uploading file"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateVersion(expectedRuntime, pathToKrt)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error uploading KRT")
}

func TestInteractor_ListRuntimes(t *testing.T) {
	s := newTestInteractorSuite(t)
	runtimeName := "test-runtime"
	expectedRuntimes := []runtime.Runtime{{Name: runtimeName}}

	s.mocks.kreClient.EXPECT().Runtime().Return(s.mocks.runtime)
	s.mocks.runtime.EXPECT().List().Return(expectedRuntimes, nil)
	s.mocks.renderer.EXPECT().RenderRuntimes(expectedRuntimes)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListRuntimes()
	assert.NoError(t, err)
}

func TestInteractor_ListRuntimes_NoRuntimesFound(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Runtime().Return(s.mocks.runtime)
	s.mocks.runtime.EXPECT().List().Return(nil, nil)
	s.mocks.renderer.EXPECT().RenderRuntimes(nil)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListRuntimes()
	assert.NoError(t, err)
}

func TestInteractor_ListRuntimes_ErrorGettingRuntimes(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Runtime().Return(s.mocks.runtime)
	s.mocks.runtime.EXPECT().List().Return(nil, errors.New("graphql error"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListRuntimes()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error getting runtimes")
}

func TestInteractor_CreateRuntime(t *testing.T) {
	s := newTestInteractorSuite(t)
	runtimeName := "test-runtime"
	description := "Test description"

	s.mocks.kreClient.EXPECT().Runtime().Return(s.mocks.runtime)
	s.mocks.runtime.EXPECT().Create(runtimeName, description).Return(nil)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateRuntime(runtimeName, description)
	assert.NoError(t, err)
}

func TestInteractor_CreateRuntime_ClientError(t *testing.T) {
	s := newTestInteractorSuite(t)
	runtimeName := "test-runtime"
	description := "Test description"

	s.mocks.kreClient.EXPECT().Runtime().Return(s.mocks.runtime)
	s.mocks.runtime.EXPECT().Create(runtimeName, description).Return(errors.New("graphql error"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateRuntime(runtimeName, description)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error creating runtime")
}
