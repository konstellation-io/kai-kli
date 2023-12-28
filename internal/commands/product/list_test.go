package product_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
)

const (
	_productID          = "productID"
	_serverName         = "serverName"
	_productDescription = "productDescription"
	_errTest            = "error test"
)

type ListProductSuite struct {
	suite.Suite

	logger               *mocks.MockLogger
	renderer             *mocks.MockRenderer
	productClient        *mocks.MockProductClient
	versionClient        *mocks.MockVersionClient
	handler              *product.Handler
	productConfiguration *productconfiguration.ProductConfigService
	kaiConfiguration     *configuration.KaiConfigService
	server               *configuration.Server

	tmpDir string
}

func TestListProductSuite(t *testing.T) {
	suite.Run(t, new(ListProductSuite))
}

func (s *ListProductSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	s.renderer = mocks.NewMockRenderer(ctrl)
	s.productClient = mocks.NewMockProductClient(s.T())
	s.versionClient = mocks.NewMockVersionClient(s.T())
	s.productConfiguration = productconfiguration.NewProductConfigService(s.logger)
	s.kaiConfiguration = configuration.NewKaiConfigService(s.logger)

	mocks.AddLoggerExpects(s.logger)

	tmpDir, err := os.MkdirTemp(s.T().TempDir(), "")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	err = os.Mkdir(path.Dir(kaiConfigPath), os.FileMode(0750))
	s.Require().NoError(err)

	_, err = os.Create(kaiConfigPath)
	s.Require().NoError(err)

	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.tmpDir = tmpDir
	s.server = &configuration.Server{Name: _serverName, Host: "http://test.com", IsDefault: true}

	kaiConf, err := s.kaiConfiguration.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(s.server)
	s.Require().NoError(err)

	err = s.kaiConfiguration.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	s.handler = product.NewHandler(
		s.logger,
		s.renderer,
		s.productClient,
		s.versionClient,
		s.productConfiguration,
	)
}

func (s *ListProductSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)

	err = os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ListProductSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ListProductSuite) TestListProducts() {
	// GIVEN
	products := []kai.Product{
		{
			ID:          _productID,
			Name:        _productName,
			Description: _productDescription,
		},
		{
			ID:          _productID,
			Name:        _productName,
			Description: _productDescription,
		},
	}
	s.productClient.EXPECT().GetProducts(s.server).Return(products, nil).Times(1)
	s.renderer.EXPECT().RenderProducts(products)

	// WHEN
	err := s.handler.ListProducts(_serverName)

	// THEN
	s.Assert().NoError(err)
}

func (s *ListProductSuite) TestListProducts_Error() {
	// GIVEN
	errTest := configuration.ErrServerNotFound
	s.productClient.EXPECT().GetProducts(s.server).Return(nil, errTest).Times(1)

	// WHEN
	err := s.handler.ListProducts(_serverName)

	// THEN
	s.Assert().EqualError(err, errTest.Error())
}
