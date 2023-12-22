package product_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type CreateProductSuite struct {
	suite.Suite

	renderer             *mocks.MockRenderer
	logger               *mocks.MockLogger
	productClient        *mocks.MockProductClient
	handler              *product.Handler
	productConfiguration *productconfiguration.ProductConfigService
	kaiConfiguration     *configuration.KaiConfigService
	server               *configuration.Server

	tmpDir string
}

func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(CreateProductSuite))
}

func (s *CreateProductSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	s.renderer = mocks.NewMockRenderer(ctrl)
	s.productClient = mocks.NewMockProductClient(s.T())
	s.kaiConfiguration = configuration.NewKaiConfigService(s.logger)
	s.productConfiguration = productconfiguration.NewProductConfigService(s.logger)

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
	s.server = &configuration.Server{Name: "test", Host: "http://test.com", IsDefault: true}

	kaiConf, err := s.kaiConfiguration.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(s.server)
	s.Require().NoError(err)

	err = s.kaiConfiguration.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	s.handler = product.NewHandler(s.logger, s.renderer, s.productClient, s.productConfiguration)
}

func (s *CreateProductSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(".kai")
	s.Require().NoError(err)

	err = os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *CreateProductSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(".kai"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *CreateProductSuite) TestCreateProduct() {
	const (
		productName = "test-product"
		version     = "v.1.0.1-test"
		description = ""
		initLocal   = false
		localPath   = ""
	)

	s.productClient.EXPECT().CreateProduct(s.server, productName, description).Return(nil).Once()
	s.logger.EXPECT().IsJSONOutputFormat().Return(false)

	err := s.handler.CreateProduct(&product.CreateProductOpts{
		Server:      s.server.Name,
		ProductName: productName,
		Version:     version,
		Description: description,
		InitLocal:   initLocal,
		LocalPath:   localPath,
	})
	s.Assert().NoError(err)
}

func (s *CreateProductSuite) TestCreateProduct_ErrorServerDoesNotExist() {
	const (
		productName = "test-product"
		version     = "v.1.0.1-test"
		description = ""
		initLocal   = false
		localPath   = ""
	)

	err := s.handler.CreateProduct(&product.CreateProductOpts{
		Server:      "inexistent server",
		ProductName: productName,
		Version:     version,
		Description: description,
		InitLocal:   initLocal,
		LocalPath:   localPath,
	})
	s.Assert().ErrorIs(err, configuration.ErrServerNotFound)
}

func (s *CreateProductSuite) TestCreateProduct_ProductClientError() {
	const (
		productName = "test-product"
		version     = "v.1.0.1-test"
		description = ""
		initLocal   = false
		localPath   = ""
	)

	expectedError := errors.New("client error")

	s.productClient.EXPECT().CreateProduct(s.server, productName, description).Return(expectedError)

	err := s.handler.CreateProduct(&product.CreateProductOpts{
		Server:      s.server.Name,
		ProductName: productName,
		Version:     version,
		Description: description,
		InitLocal:   initLocal,
		LocalPath:   localPath,
	})
	s.Assert().ErrorIs(err, expectedError)
}

func (s *CreateProductSuite) TestCreateProduct_InitLocalDefaultPath() {
	const (
		productName = "test-product"
		version     = "v.1.0.1-test"
		description = ""
		initLocal   = true
	)

	tmpProductPath, err := os.MkdirTemp(s.T().TempDir(), "")
	s.Require().NoError(err)

	defer os.RemoveAll(tmpProductPath)

	s.productClient.EXPECT().CreateProduct(s.server, productName, description).Return(nil).Once()
	s.logger.EXPECT().IsJSONOutputFormat().Return(false)

	err = s.handler.CreateProduct(&product.CreateProductOpts{
		Server:      s.server.Name,
		ProductName: productName,
		Version:     version,
		Description: description,
		InitLocal:   initLocal,
		LocalPath:   tmpProductPath,
	})
	s.Assert().NoError(err)

	cfg, err := s.productConfiguration.GetConfiguration(productName, tmpProductPath)
	s.Require().NoError(err)

	s.Assert().Equal(&productconfiguration.KaiProductConfiguration{
		Krt: &krt.Krt{
			Version:     version,
			Description: description,
			Config:      map[string]string{},
			Workflows:   []krt.Workflow{},
		}}, cfg)
}

func (s *CreateProductSuite) TestCreateProduct_LocalConfigAlreadyExists() {
	const (
		productName = "test-product"
		version     = "v.1.0.1-test"
		description = ""
		initLocal   = true
	)

	tmpProductPath, err := os.MkdirTemp(s.T().TempDir(), "")
	s.Require().NoError(err)

	defer os.RemoveAll(tmpProductPath)

	err = s.productConfiguration.WriteConfiguration(&productconfiguration.KaiProductConfiguration{
		Krt: &krt.Krt{
			Version: version,
		},
	}, productName, tmpProductPath)
	s.Require().NoError(err)

	s.productClient.EXPECT().CreateProduct(s.server, productName, description).Return(nil).Once()
	s.logger.EXPECT().IsJSONOutputFormat().Return(false)

	err = s.handler.CreateProduct(&product.CreateProductOpts{
		Server:      s.server.Name,
		ProductName: productName,
		Version:     version,
		Description: description,
		InitLocal:   initLocal,
		LocalPath:   tmpProductPath,
	})
	s.Require().Error(err, product.ErrProductConfigurationAlreadyExists)
}
