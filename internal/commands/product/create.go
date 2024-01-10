package product

import (
	"errors"
	"fmt"

	"github.com/konstellation-io/kli/internal/services/configuration"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/krt/pkg/krt"
)

var ErrProductConfigurationAlreadyExists = errors.New("product configuration already exists")

type CreateProductOpts struct {
	Server      string
	ProductName string
	Version     string
	Description string
	InitLocal   bool
	LocalPath   string
}

func (h *Handler) CreateProduct(opts *CreateProductOpts) error {
	configService := configuration.NewKaiConfigService(h.logger)

	conf, err := configService.GetConfiguration()
	if err != nil {
		return err
	}

	server, err := conf.GetServerOrDefault(opts.Server)
	if err != nil {
		return err
	}

	err = h.productClient.CreateProduct(server, opts.ProductName, opts.Description)
	if err != nil {
		return fmt.Errorf("creating remote product: %w", err)
	}

	if h.logger.IsJSONOutputFormat() {
		h.renderer.RenderProductCreated(opts.ProductName, server, opts.InitLocal)
	} else {
		h.logger.Success(fmt.Sprintf("Product %s successfully created in the server %s", opts.ProductName, server.Name))

		if opts.InitLocal {
			if err := h.createLocalProductConfiguration(opts); err != nil {
				return fmt.Errorf("creating local product's configuration: %w", err)
			}

			h.logger.Success(fmt.Sprintf("Product %q local configuration created.", opts.ProductName))
		}
	}

	return nil
}

func (h *Handler) createLocalProductConfiguration(opts *CreateProductOpts) error {
	cfg, err := h.configService.GetConfiguration(opts.ProductName, opts.LocalPath)
	if err != nil && !errors.Is(err, productconfiguration.ErrProductConfigNotFound) {
		return err
	}

	if cfg != nil {
		return ErrProductConfigurationAlreadyExists
	}

	err = h.configService.WriteConfiguration(
		&productconfiguration.KaiProductConfiguration{
			Krt: &krt.Krt{
				Version:     opts.Version,
				Description: opts.Description,
			},
		}, opts.ProductName, opts.LocalPath)
	if err != nil {
		return err
	}

	return nil
}
