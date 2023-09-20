package version_test

import (
	"errors"
	"time"

	apiVersion "github.com/konstellation-io/kli/api/version"
	"github.com/konstellation-io/kli/internal/commands/version"
)

func (s *VersionSuite) TestGetVersion() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
	)

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
	productName := "test-product"
	versionTag := "v.1.0.1-test"
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Get(s.server, productName, versionTag).Return(nil, expectedError).Once()

	err := s.handler.GetVersion(&version.GetVersionOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
	})
	s.Assert().EqualError(err, expectedError.Error())
}
