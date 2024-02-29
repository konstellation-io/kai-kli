//go:build unit

package product_test

import (
	"errors"
	"os"

	"github.com/konstellation-io/kli/internal/commands/product"
	"github.com/konstellation-io/kli/internal/entity"
)

const (
	_productName = "test-product"
	_versionTag  = "v1.0.0"
)

func (s *ProductSuite) TestBindProduct() {
	testVersion := &entity.Version{
		Tag: _versionTag,
	}

	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      false,
	}

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(testVersion, nil).Once()
	s.renderer.EXPECT().RenderProductBinded(_productName).Times(1)

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
	testVersion := &entity.Version{
		Tag: _versionTag,
	}

	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      false,
	}

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(testVersion, nil).Times(1)
	s.renderer.EXPECT().RenderProductBinded(_productName).Times(1)

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
	testVersion := &entity.Version{
		Tag: _versionTag,
	}

	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      true,
	}

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(testVersion, nil).Times(2)
	s.renderer.EXPECT().RenderProductBinded(_productName).Times(2)

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

func (s *ProductSuite) TestBind_InvalidServer() {
	bindOpts := &product.BindProductOpts{
		Server:     "invalid-server",
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      true,
	}

	err := s.handler.Bind(bindOpts)
	s.Error(err)
}

func (s *ProductSuite) TestBind_GetVersionError() {
	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil,
		LocalPath:  "",
		Force:      true,
	}

	testErr := errors.New("test error")

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(nil, testErr).Times(1)

	err := s.handler.Bind(bindOpts)
	s.ErrorIs(err, testErr)
}

func (s *ProductSuite) TestBind_GetVersionError_ProductExists() {
	bindOpts := &product.BindProductOpts{
		Server:     s.server.Name,
		ProductID:  _productName,
		VersionTag: nil, // GIVEN no specific version tag
		LocalPath:  "",
		Force:      true,
	}

	s.versionClient.EXPECT().Get(s.server, _productName, bindOpts.VersionTag).Return(nil, product.ErrProductExistsNoVersions).Times(1)
	s.renderer.EXPECT().RenderProductBinded(_productName).Times(1)

	err := s.handler.Bind(bindOpts)
	s.NoError(err)
}
