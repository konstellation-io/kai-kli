package version_test

import (
	"errors"

	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/entity"
)

func (s *VersionSuite) TestPublishVersion() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
	)

	urls := []entity.TriggerEndpoint{
		{
			Trigger: "test-trigger",
			URL:     "test-url",
		},
	}

	s.versionClient.EXPECT().Publish(s.server, productName, versionTag, comment).Return(urls, nil).Once()
	s.renderer.EXPECT().RenderTriggers(urls)

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestPublishVersion_NoTriggersPublished() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
	)

	s.versionClient.EXPECT().Publish(s.server, productName, versionTag, comment).Return(nil, nil).Once()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestPublishVersion_ClientError() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
	)

	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Publish(s.server, productName, versionTag, comment).Return(nil, expectedError).Once()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().ErrorIs(err, expectedError)
}
