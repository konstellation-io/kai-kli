//go:build unit

package version_test

import (
	"errors"
	"time"

	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/entity"
)

func (s *VersionSuite) TestListVersion_NoFilter_OK() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
	)

	versions := []*entity.Version{
		{
			Tag:          versionTag,
			CreationDate: time.Now(),
			Status:       "CREATED",
		},
	}

	s.versionClient.EXPECT().List(s.server, productName, (*string)(nil)).Return(versions, nil).Once()
	s.renderer.EXPECT().RenderVersions(productName, versions)

	err := s.handler.ListVersions(&version.ListVersionsOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestListVersion_WithFilter_Ok() {
	const (
		productName = "test-product"
	)

	versions := []*entity.Version{
		{
			Tag:          "v1.0.0",
			CreationDate: time.Now(),
			Status:       entity.VersionStatusStopped,
		},
		{
			Tag:          "v1.0.1",
			CreationDate: time.Now(),
			Status:       entity.VersionStatusStarted,
		},
	}

	statusFilter := entity.VersionStatusStarted
	listVersionOpts := &version.ListVersionsOpts{
		ServerName:   s.server.Name,
		ProductID:    productName,
		StatusFilter: &statusFilter,
	}

	filteredVersions := []*entity.Version{versions[1]}

	s.versionClient.EXPECT().List(s.server, productName, listVersionOpts.StatusFilter).Return(filteredVersions, nil).Once()
	s.renderer.EXPECT().RenderVersions(productName, filteredVersions)

	err := s.handler.ListVersions(listVersionOpts)
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestListVersion_WrongServer_ExpectError() {
	err := s.handler.ListVersions(&version.ListVersionsOpts{
		ServerName: "invalid-server",
		ProductID:  productName,
	})
	s.Assert().Error(err)
}

func (s *VersionSuite) TestListVersion_ClientFails_ExpectErrors() {
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().List(s.server, productName, (*string)(nil)).Return(nil, expectedError).Once()

	err := s.handler.ListVersions(&version.ListVersionsOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
	})
	s.Assert().EqualError(err, expectedError.Error())
}
