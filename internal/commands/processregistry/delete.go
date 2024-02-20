package processregistry

import (
	"errors"
	"fmt"

	"github.com/konstellation-io/kli/internal/services/configuration"
)

var (
	ErrCannotbePublicAndProduct = errors.New("process cannot be public and have a product")
	ErrNoProduct                = errors.New("product flag is required for product processes")
)

type DeleteProcessOpts struct {
	ServerName string
	ProductID  string
	ProcessID  string
	Version    string
	IsPublic   bool
}

func (opts *DeleteProcessOpts) Validate() error {
	if opts.IsPublic && opts.ProductID != "" {
		return ErrCannotbePublicAndProduct
	}

	if !opts.IsPublic && opts.ProductID == "" {
		return ErrNoProduct
	}

	return nil
}

func (c *Handler) DeleteProcess(opts *DeleteProcessOpts) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validating args: %w", err)
	}

	kaiConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	deletedProcess, err := c.deleteProcess(srv, opts)
	if err != nil {
		return err
	}

	c.renderer.RenderProcessDeleted(deletedProcess)

	return nil
}

func (c *Handler) deleteProcess(
	srv *configuration.Server, opts *DeleteProcessOpts,
) (string, error) {
	if opts.IsPublic {
		return c.processRegistryClient.DeletePublic(srv, opts.ProcessID, opts.Version)
	}

	return c.processRegistryClient.Delete(srv, opts.ProductID, opts.ProcessID, opts.Version)
}
