package version_test

import (
	"errors"

	"github.com/konstellation-io/kli/internal/commands/version"
)

func (s *VersionSuite) TestStopVersion() {
	s.versionClient.EXPECT().Stop(s.server, productName, versionTag, comment).Return(versionTag, nil).Once()

	err := s.handler.Stop(&version.StopOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestStopVersion_ErrorIfClientFails() {
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Stop(s.server, productName, versionTag, comment).Return("", expectedError).Once()

	err := s.handler.Stop(&version.StopOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().EqualError(err, expectedError.Error())
}
