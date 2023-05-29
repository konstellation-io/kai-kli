package kre_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kli/api/kre/product"
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
	product   *mocks.MockProductInterface
	renderer  *mocks.MockRenderer
}

func newTestInteractorSuite(t *testing.T) *testInteractorSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	kreClient := mocks.NewMockKreInterface(ctrl)
	versionMock := mocks.NewMockVersionInterface(ctrl)
	productMock := mocks.NewMockProductInterface(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)

	return &testInteractorSuite{
		ctrl: ctrl,
		mocks: suiteMocks{
			logger:    logger,
			kreClient: kreClient,
			version:   versionMock,
			product:   productMock,
			renderer:  renderer,
		},
	}
}

func TestInteractor_ListVersions(t *testing.T) {
	s := newTestInteractorSuite(t)

	versionProduct := version.Version{
		Name:   "test-version",
		Status: "STARTED",
	}
	versionList := version.List([]version.Version{versionProduct})

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-product").Return(versionList, nil)
	s.mocks.renderer.EXPECT().RenderVersions(versionList).Return()

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-product")
	assert.NoError(t, err)
}

func TestInteractor_ListVersions_NoVersionsInProduct(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-product").Return(nil, nil)
	s.mocks.renderer.EXPECT().RenderVersions(nil).Return()

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-product")
	assert.NoError(t, err)
}

func TestInteractor_ListVersions_RequestFail(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-product").Return(nil, errors.New("client error"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-product")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error listing versions")
}

func TestInteractor_CreateVersion(t *testing.T) {
	s := newTestInteractorSuite(t)
	pathToKrt := "path/to/krt"
	expectedProduct := "test-product"

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().Create(expectedProduct, pathToKrt).Return("version-name", nil)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateVersion(expectedProduct, pathToKrt)
	assert.NoError(t, err)
}

func TestInteractor_CreateVersion_ClientError(t *testing.T) {
	s := newTestInteractorSuite(t)
	pathToKrt := "path/to/krt"
	expectedProduct := "test-product"

	s.mocks.kreClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().Create(expectedProduct, pathToKrt).Return("", errors.New("error uploading file"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateVersion(expectedProduct, pathToKrt)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error uploading KRT")
}

func TestInteractor_ListProducts(t *testing.T) {
	s := newTestInteractorSuite(t)
	productName := "test-product"
	expectedProducts := []product.Product{{Name: productName}}

	s.mocks.kreClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().List().Return(expectedProducts, nil)
	s.mocks.renderer.EXPECT().RenderProducts(expectedProducts)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListProducts()
	assert.NoError(t, err)
}

func TestInteractor_ListProducts_NoProductsFound(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().List().Return(nil, nil)
	s.mocks.renderer.EXPECT().RenderProducts(nil)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListProducts()
	assert.NoError(t, err)
}

func TestInteractor_ListProducts_ErrorGettingProducts(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kreClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().List().Return(nil, errors.New("graphql error"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.ListProducts()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error getting products")
}

func TestInteractor_CreateProduct(t *testing.T) {
	s := newTestInteractorSuite(t)
	productName := "test-product"
	description := "Test description"

	s.mocks.kreClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().Create(productName, description).Return(nil)

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateProduct(productName, description)
	assert.NoError(t, err)
}

func TestInteractor_CreateProduct_ClientError(t *testing.T) {
	s := newTestInteractorSuite(t)
	productName := "test-product"
	description := "Test description"

	s.mocks.kreClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().Create(productName, description).Return(errors.New("graphql error"))

	krtInteractor := kre.NewInteractor(s.mocks.logger, s.mocks.kreClient, s.mocks.renderer)
	err := krtInteractor.CreateProduct(productName, description)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error creating product")
}
