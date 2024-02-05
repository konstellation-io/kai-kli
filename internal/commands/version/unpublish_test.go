//go:build unit

package version_test

import (
	"errors"

	"github.com/konstellation-io/kli/internal/commands/version"
)

func (s *VersionSuite) TestUnpublishVersion() {
	s.versionClient.EXPECT().Unpublish(s.server, productName, versionTag, comment).Return(versionTag, nil).Once()
	s.renderer.EXPECT().RenderUnpublishVersion(productName, versionTag).Times(1)

	err := s.handler.Unpublish(&version.UnpublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestUnpublishVersion_ErrorIfClientFails() {
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Unpublish(s.server, productName, versionTag, comment).Return("", expectedError).Once()

	err := s.handler.Unpublish(&version.UnpublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().EqualError(err, expectedError.Error())
}
