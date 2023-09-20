package version_test

import (
	"errors"
	"time"

	apiVersion "github.com/konstellation-io/kli/api/version"
	"github.com/konstellation-io/kli/internal/commands/version"
)

func (s *VersionSuite) TestListVersion() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
	)

	versions := []*apiVersion.Version{
		{
			Tag:          versionTag,
			CreationDate: time.Now(),
			Status:       "CREATED",
		},
	}

	s.versionClient.EXPECT().List(s.server, productName).Return(versions, nil).Once()
	s.renderer.EXPECT().RenderVersions(versions)

	err := s.handler.ListVersions(&version.ListVersionsOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestListVersion_ErrorIfClientFails() {
	productName := "test-product"
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().List(s.server, productName).Return(nil, expectedError).Once()

	err := s.handler.ListVersions(&version.ListVersionsOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
	})
	s.Assert().EqualError(err, expectedError.Error())
}
