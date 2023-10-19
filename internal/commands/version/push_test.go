package version_test

import (
	"errors"
	"path"

	"github.com/konstellation-io/kli/internal/commands/version"
)

func (s *VersionSuite) TestCreateProduct() {
	krtFilePath := path.Join("testdata", productName+".yaml")

	s.versionClient.EXPECT().Push(s.server, productName, krtFilePath).Return(versionTag, nil).Once()

	err := s.handler.PushVersion(&version.PushVersionOpts{
		Server:      s.server.Name,
		KrtFilePath: krtFilePath,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestCreateProduct_ErrorIfClientFails() {
	krtFilePath := path.Join("testdata", productName+".yaml")
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Push(s.server, productName, krtFilePath).Return("", expectedError).Once()

	err := s.handler.PushVersion(&version.PushVersionOpts{
		Server:      s.server.Name,
		KrtFilePath: krtFilePath,
	})
	s.Assert().ErrorIs(err, expectedError)
}
