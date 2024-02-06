//go:build unit

package version_test

import (
	"errors"

	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/entity"
)

func (s *VersionSuite) TestPublishVersion() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusStarted,
		}
		product = &kai.Product{PublishedVersion: nil}
		urls    = []entity.TriggerEndpoint{
			{
				Trigger: "test-trigger",
				URL:     "test-url",
			},
		}
	)

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()
	s.productClient.EXPECT().GetProduct(s.server, productName).Return(product, nil).Once()
	s.versionClient.EXPECT().Publish(s.server, productName, versionTag, comment, false).Return(urls, nil).Once()
	s.renderer.EXPECT().RenderTriggers(urls).Return()
	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestPublishVersion_NoTriggersPublished() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusStarted,
		}
		product = &kai.Product{PublishedVersion: nil}
	)

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()
	s.productClient.EXPECT().GetProduct(s.server, productName).Return(product, nil).Once()
	s.versionClient.EXPECT().Publish(s.server, productName, versionTag, comment, false).Return(nil, nil).Once()
	s.renderer.EXPECT().RenderTriggers(nil).Return()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestPublishVersion_ClientError() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusStarted,
		}
		product = &kai.Product{PublishedVersion: nil}
	)

	expectedError := errors.New("client error")

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()
	s.productClient.EXPECT().GetProduct(s.server, productName).Return(product, nil).Once()
	s.versionClient.EXPECT().Publish(s.server, productName, versionTag, comment, false).Return(nil, expectedError).Once()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().ErrorIs(err, expectedError)
}

func (s *VersionSuite) TestPublishVersion_VersionNotStarted() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusStopped,
		}
	)

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().ErrorIs(err, version.ErrVersionNotStarted)
}

func (s *VersionSuite) TestPublishVersion_VersionNotStarted_Force() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusStopped,
		}
	)

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
		Force:      true,
	})
	s.Assert().ErrorIs(err, version.ErrVersionNotStarted)
}

func (s *VersionSuite) TestPublishVersion_VersionAlreadyPublished() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusPublished,
		}
	)

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
		Force:      true,
	})
	s.Assert().ErrorIs(err, version.ErrVersionAlreadyPublished)
}

func (s *VersionSuite) TestPublishVersion_ProductAlreadyPublished_NoForce() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		oldVersion  = "old-version"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusStarted,
		}
		product = &kai.Product{PublishedVersion: &oldVersion}
	)

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()
	s.productClient.EXPECT().GetProduct(s.server, productName).Return(product, nil).Once()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
	})
	s.Assert().ErrorIs(err, version.ErrProductAlreadyPublished)
}

func (s *VersionSuite) TestPublishVersion_ProductAlreadyPublished_Force() {
	var (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
		oldVersion  = "old-version"
		comment     = "Test comment"
		vers        = &entity.Version{
			Tag:    versionTag,
			Status: entity.VersionStatusStarted,
		}
		product = &kai.Product{PublishedVersion: &oldVersion}
	)

	s.versionClient.EXPECT().Get(s.server, productName, &versionTag).Return(vers, nil).Once()
	s.productClient.EXPECT().GetProduct(s.server, productName).Return(product, nil).Once()
	s.versionClient.EXPECT().Publish(s.server, productName, versionTag, comment, true).Return(nil, nil).Once()
	s.renderer.EXPECT().RenderTriggers(nil).Return()

	err := s.handler.Publish(&version.PublishOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
		VersionTag: versionTag,
		Comment:    comment,
		Force:      true,
	})
	s.Assert().NoError(err)
}
