package version_test

import (
	"errors"

	"github.com/konstellation-io/kli/internal/commands/version"
)

func (s *VersionSuite) TestStartVersion() {
	s.versionClient.EXPECT().Start(s.server, productName, versionTag, comment).Return(versionTag, nil).Once()

	err := s.handler.Start(&version.StartOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestStartVersion_ErrorIfClientFails() {
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Start(s.server, productName, versionTag, comment).Return("", expectedError).Once()

	err := s.handler.Start(&version.StartOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().EqualError(err, expectedError.Error())
}
