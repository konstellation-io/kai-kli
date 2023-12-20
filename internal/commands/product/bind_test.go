package product_test

import (
	"os"

	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/entity"
)

const _productName = "test-product"

var _testVersion = &entity.Version{
	Tag: "v.1.0.0",
}

func (s *ProductSuite) TestBindProduct() {
	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      false,
	}

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(_testVersion, nil).Once()

	err := s.handler.Bind(bindOpts)
	s.Require().NoError(err)

	fileName := _productName + ".yaml"

	// Cleanup file
	_, err = os.Stat(fileName)
	s.Require().NoError(err)

	err = os.Remove(fileName)
	s.Require().NoError(err)
}

func (s *ProductSuite) TestBindProductAlreadyExists() {
	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      false,
	}

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(_testVersion, nil)

	err := s.handler.Bind(bindOpts)
	s.Require().NoError(err)

	err = s.handler.Bind(bindOpts)
	s.Assert().Error(err)

	fileName := _productName + ".yaml"

	// Cleanup file
	_, err = os.Stat(fileName)
	s.Require().NoError(err)

	err = os.Remove(fileName)
	s.Require().NoError(err)
}

func (s *ProductSuite) TestBindForce() {
	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      true,
	}

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(_testVersion, nil)

	err := s.handler.Bind(bindOpts)
	s.Require().NoError(err)

	err = s.handler.Bind(bindOpts)
	s.Require().NoError(err)

	fileName := _productName + ".yaml"

	// Cleanup file
	_, err = os.Stat(fileName)
	s.Require().NoError(err)

	err = os.Remove(fileName)
	s.Require().NoError(err)
}
