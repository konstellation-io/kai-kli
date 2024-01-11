//go:build unit

package version_test

import (
	"errors"

	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/entity"
)

func (s *VersionSuite) TestStartVersion() {
	startOpts := &version.StartOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	}

	s.versionClient.EXPECT().Get(s.server, productName, &startOpts.VersionTag).Return(&entity.Version{}, nil).Once()
	s.versionClient.EXPECT().Start(s.server, productName, versionTag, comment).Return(versionTag, nil).Once()

	err := s.handler.Start(startOpts)
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestStartVersion_ErrorIfClientFails() {
	expectedError := errors.New("client error")

	startOpts := &version.StartOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	}

	s.versionClient.EXPECT().Get(s.server, productName, &startOpts.VersionTag).Return(&entity.Version{}, nil).Once()
	s.versionClient.EXPECT().Start(s.server, productName, versionTag, comment).Return("", expectedError).Once()

	err := s.handler.Start(startOpts)
	s.Assert().EqualError(err, expectedError.Error())
}
