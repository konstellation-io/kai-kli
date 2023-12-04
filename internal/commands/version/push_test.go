package version_test

import (
	"errors"
	"github.com/konstellation-io/kli/internal/commands/version"
)

func (s *VersionSuite) TestPushVersion_ExpectOK() {
	s.versionClient.EXPECT().Push(s.server, productName, s.krtFilePath).Return(versionTag, nil).Once()

	err := s.handler.PushVersion(&version.PushVersionOpts{
		Server:      s.server.Name,
		KrtFilePath: s.krtFilePath,
	})

	s.Assert().NoError(err)
}

func (s *VersionSuite) TestPushVersion_UpdateVersionOnPush_ExpectOK() {
	defaultProductVersion := "v1.0.0"
	newVersion := "v2.0.0"
	newDescription := "New description"

	s.versionClient.EXPECT().Push(s.server, productName, s.krtFilePath).Return(newVersion, nil).Once()

	err := s.handler.PushVersion(&version.PushVersionOpts{
		Server:      s.server.Name,
		KrtFilePath: s.krtFilePath,
		Version:     newVersion,
		Description: newDescription,
	})

	s.Assert().NoError(err)

	productConfig, err := s.productConfig.GetConfiguration(productName, "testdata")
	s.Require().NoError(err)

	productVersion := productConfig.GetProductVersion()
	productDescription := productConfig.GetProductDescription()

	s.Assert().Equal(newVersion, productVersion)
	s.Assert().Equal(newDescription, productDescription)

	err = productConfig.UpdateProductVersion(defaultProductVersion)
	s.Require().NoError(err)
}

func (s *VersionSuite) TestPushVersion_ErrorIfClientFails() {
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Push(s.server, productName, s.krtFilePath).Return("", expectedError).Once()

	err := s.handler.PushVersion(&version.PushVersionOpts{
		Server:      s.server.Name,
		KrtFilePath: s.krtFilePath,
	})
	s.Assert().ErrorIs(err, expectedError)
}
