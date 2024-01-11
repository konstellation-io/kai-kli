//go:build unit

package product_test

import (
	"errors"
	"os"

	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/krt/pkg/krt"
)

func (s *ProductSuite) TestCreateProduct() {
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

func (s *ProductSuite) TestCreateProduct_ErrorServerDoesNotExist() {
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

func (s *ProductSuite) TestCreateProduct_ProductClientError() {
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

func (s *ProductSuite) TestCreateProduct_InitLocalDefaultPath() {
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

func (s *ProductSuite) TestCreateProduct_LocalConfigAlreadyExists() {
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
