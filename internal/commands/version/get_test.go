package version_test

import (
	"errors"
	"time"

	apiVersion "github.com/konstellation-io/kli/api/version"
	"github.com/konstellation-io/kli/internal/commands/version"
)

const (
	productName = "test-product"
	versionTag  = "v.1.0.1-test"
)

func (s *VersionSuite) TestGetVersion() {
	oneVersion := &apiVersion.Version{
		Tag:          versionTag,
		CreationDate: time.Now(),
		Status:       "CREATED",
	}

	s.versionClient.EXPECT().Get(s.server, productName, versionTag).Return(oneVersion, nil).Once()
	s.renderer.EXPECT().RenderVersion(productName, oneVersion)

	err := s.handler.GetVersion(&version.GetVersionOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestGetVersion_ErrorIfClientFails() {
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Get(s.server, productName, versionTag).Return(nil, expectedError).Once()

	err := s.handler.GetVersion(&version.GetVersionOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
	})
	s.Assert().EqualError(err, expectedError.Error())
}

func (s *VersionSuite) TestGetVersion_ContainsError() {
	oneVersion := &apiVersion.Version{
		Tag:          versionTag,
		CreationDate: time.Now(),
		Status:       "CREATED",
		Error:        "error",
	}

	s.versionClient.EXPECT().Get(s.server, productName, versionTag).Return(oneVersion, nil).Once()
	s.renderer.EXPECT().RenderVersion(productName, oneVersion)

	err := s.handler.GetVersion(&version.GetVersionOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
	})
	s.Assert().NoError(err)
}
