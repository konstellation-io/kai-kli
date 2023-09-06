package product

import (
	"errors"

	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/krt/pkg/krt"
)

type BindProductOpts struct {
	Server    string
	ProductID string
	LocalPath string
	Force     bool
}

func (h *Handler) Bind(opts *BindProductOpts) error {
	// Get remote product
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(opts.Server)
	if err != nil {
		return err
	}

	product, err := h.productClient.GetProduct(server, opts.ProductID)
	if err != nil {
		return err
	}

	// Check if local config already exists
	// - if exists and --force, override
	// - if not exists, create a new config
	productConfig, err := h.configService.GetConfiguration(product.ID, opts.LocalPath)
	if err != nil && !errors.Is(err, productconfiguration.ErrProductConfigNotFound) {
		return err
	}

	if productConfig != nil && !opts.Force {
		return ErrProductConfigurationAlreadyExists
	}

	err = h.configService.WriteConfiguration(&productconfiguration.KaiProductConfiguration{
		Krt: &krt.Krt{
			Description: product.Description,
		},
	}, product.ID, opts.LocalPath)

	if err != nil {
		return err
	}

	h.logger.Success("Product successfully binded!")

	return nil
}
