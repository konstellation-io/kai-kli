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

	getOpts := &version.GetVersionOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
	}

	s.versionClient.EXPECT().Get(s.server, productName, &getOpts.VersionTag).Return(testVersion, nil).Once()
	s.renderer.EXPECT().RenderVersion(productName, testVersion)

	err := s.handler.GetVersion(getOpts)
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestGetVersion_ErrorIfClientFails() {
	expectedError := errors.New("client error")

	getOpts := &version.GetVersionOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
	}

	s.versionClient.EXPECT().Get(s.server, productName, &getOpts.VersionTag).Return(nil, expectedError).Once()

	err := s.handler.GetVersion(getOpts)
	s.Assert().EqualError(err, expectedError.Error())
}

func (s *VersionSuite) TestGetVersion_ContainsError() {
	oneVersion := &entity.Version{
		Tag:          versionTag,
		CreationDate: time.Now(),
		Status:       "CREATED",
		Error:        "error",
	}

	getOpts := &version.GetVersionOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
	}

	s.versionClient.EXPECT().Get(s.server, productName, &getOpts.VersionTag).Return(oneVersion, nil).Once()
	s.renderer.EXPECT().RenderVersion(productName, oneVersion)

	err := s.handler.GetVersion(getOpts)
	s.Assert().NoError(err)
}
