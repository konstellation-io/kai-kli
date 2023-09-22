package version_test

import (
	"errors"
	"time"

	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/entity"
)

func (s *VersionSuite) TestGetVersion() {
	testVersion := &entity.Version{
		Tag:          versionTag,
		CreationDate: time.Now(),
		Status:       "CREATED",
	}

	s.versionClient.EXPECT().Get(s.server, productName, versionTag).Return(testVersion, nil).Once()
	s.renderer.EXPECT().RenderVersion(productName, testVersion)

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
	oneVersion := &entity.Version{
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
