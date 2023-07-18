package kai_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kli/api/kai/product"
	"github.com/konstellation-io/kli/api/kai/version"
	"github.com/konstellation-io/kli/internal/commands/kai"
	"github.com/konstellation-io/kli/mocks"
)

type testInteractorSuite struct {
	ctrl  *gomock.Controller
	mocks suiteMocks
}

type suiteMocks struct {
	logger    *mocks.MockLogger
	kaiClient *mocks.MockKaiInterface
	version   *mocks.MockVersionInterface
	product   *mocks.MockProductInterface
	renderer  *mocks.MockRenderer
}

func newTestInteractorSuite(t *testing.T) *testInteractorSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	kaiClient := mocks.NewMockKaiInterface(ctrl)
	versionMock := mocks.NewMockVersionInterface(ctrl)
	productMock := mocks.NewMockProductInterface(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)

	return &testInteractorSuite{
		ctrl: ctrl,
		mocks: suiteMocks{
			logger:    logger,
			kaiClient: kaiClient,
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

	s.mocks.kaiClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-product").Return(versionList, nil)
	s.mocks.renderer.EXPECT().RenderVersions(versionList).Return()

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-product")
	assert.NoError(t, err)
}

func TestInteractor_ListVersions_NoVersionsInProduct(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kaiClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-product").Return(nil, nil)
	s.mocks.renderer.EXPECT().RenderVersions(nil).Return()

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-product")
	assert.NoError(t, err)
}

func TestInteractor_ListVersions_RequestFail(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kaiClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().List("test-product").Return(nil, errors.New("client error"))

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.ListVersions("test-product")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error listing versions")
}

func TestInteractor_CreateVersion(t *testing.T) {
	s := newTestInteractorSuite(t)
	pathToKrt := "path/to/krt"
	expectedProduct := "test-product"

	s.mocks.kaiClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().Create(expectedProduct, pathToKrt).Return("version-name", nil)

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.CreateVersion(expectedProduct, pathToKrt)
	assert.NoError(t, err)
}

func TestInteractor_CreateVersion_ClientError(t *testing.T) {
	s := newTestInteractorSuite(t)
	pathToKrt := "path/to/krt"
	expectedProduct := "test-product"

	s.mocks.kaiClient.EXPECT().Version().Return(s.mocks.version)
	s.mocks.version.EXPECT().Create(expectedProduct, pathToKrt).Return("", errors.New("error uploading file"))

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.CreateVersion(expectedProduct, pathToKrt)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error uploading KRT")
}

func TestInteractor_ListProducts(t *testing.T) {
	s := newTestInteractorSuite(t)
	productName := "test-product"
	expectedProducts := []product.Product{{Name: productName}}

	s.mocks.kaiClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().List().Return(expectedProducts, nil)
	s.mocks.renderer.EXPECT().RenderProducts(expectedProducts)

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.ListProducts()
	assert.NoError(t, err)
}

func TestInteractor_ListProducts_NoProductsFound(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kaiClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().List().Return(nil, nil)
	s.mocks.renderer.EXPECT().RenderProducts(nil)

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.ListProducts()
	assert.NoError(t, err)
}

func TestInteractor_ListProducts_ErrorGettingProducts(t *testing.T) {
	s := newTestInteractorSuite(t)

	s.mocks.kaiClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().List().Return(nil, errors.New("graphql error"))

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.ListProducts()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error getting products")
}

func TestInteractor_CreateProduct(t *testing.T) {
	s := newTestInteractorSuite(t)
	productName := "test-product"
	description := "Test description"

	s.mocks.kaiClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().Create(productName, description).Return(nil)

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.CreateProduct(productName, description)
	assert.NoError(t, err)
}

func TestInteractor_CreateProduct_ClientError(t *testing.T) {
	s := newTestInteractorSuite(t)
	productName := "test-product"
	description := "Test description"

	s.mocks.kaiClient.EXPECT().Product().Return(s.mocks.product)
	s.mocks.product.EXPECT().Create(productName, description).Return(errors.New("graphql error"))

	krtInteractor := kai.NewInteractor(s.mocks.logger, s.mocks.kaiClient, s.mocks.renderer)
	err := krtInteractor.CreateProduct(productName, description)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error creating product")
}
