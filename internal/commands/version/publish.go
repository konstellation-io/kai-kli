package version

import (
	"errors"

	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

var (
	ErrVersionNotStarted       = errors.New("version must be started before publish")
	ErrVersionAlreadyPublished = errors.New("version already published")
	ErrProductAlreadyPublished = errors.New("product has a published version, use the --force param to replace it")
)

type PublishOpts struct {
	ServerName string
	ProductID  string
	VersionTag string
	Comment    string
	Force      bool
}

func (h *Handler) Publish(opts *PublishOpts) error {
	kaiConfig, err := h.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	err = h.validateOpts(opts, srv)
	if err != nil {
		return err
	}

	publishedTriggers, err := h.versionClient.Publish(
		srv, opts.ProductID, opts.VersionTag, opts.Comment, opts.Force,
	)
	if err != nil {
		return err
	}

	h.renderer.RenderPublishVersion(opts.ProductID, opts.VersionTag, publishedTriggers)

	return nil
}

func (h *Handler) validateOpts(opts *PublishOpts, srv *configuration.Server) error {
	err := h.checkVersionCanBeStarted(srv, opts)
	if err != nil {
		return err
	}

	return h.checkIfReplacePublishedVersion(srv, opts)
}

func (h *Handler) checkVersionCanBeStarted(srv *configuration.Server, opts *PublishOpts) error {
	versionToPublish, err := h.versionClient.Get(srv, opts.ProductID, &opts.VersionTag)
	if err != nil {
		return err
	}

	switch versionToPublish.Status {
	case entity.VersionStatusStarted:
		return nil
	case entity.VersionStatusPublished:
		return ErrVersionAlreadyPublished
	default:
		return ErrVersionNotStarted
	}
}

func (h *Handler) checkIfReplacePublishedVersion(srv *configuration.Server, opts *PublishOpts) error {
	p, err := h.productClient.GetProduct(srv, opts.ProductID)
	if err != nil {
		return err
	}

	if p.PublishedVersion != nil {
		if !opts.Force {
			return ErrProductAlreadyPublished
		}

		h.logger.Info("Published version will be replaced, remember to stop the old version free up resources.")
	}

	return nil
}
