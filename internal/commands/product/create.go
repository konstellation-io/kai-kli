package product

import (
	"errors"
	"fmt"
	"os"

	"github.com/konstellation-io/krt/pkg/krt"

	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
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

func (c *Handler) CreateProduct(opts *CreateProductOpts) error {
	productConfigPath, err := productconfiguration.GetProductConfigFilePath(opts.ProductName)
	if err != nil {
		return err
	}

	// TODO create product in KAI server

	_, err = os.Stat(productConfigPath)
	if !os.IsNotExist(err) {
		return ErrProductConfigurationAlreadyExists
	}

	err = c.configService.WriteConfiguration(
		&productconfiguration.KaiProductConfiguration{
			Krt: &krt.Krt{
				Version:     opts.ProductName,
				Description: opts.Description,
				Config:      map[string]string{},
				Workflows:   []krt.Workflow{},
			},
		}, opts.ProductName)
	if err != nil {
		return err
	}

	// TODO add renderer method to show products

	c.logger.Success(fmt.Sprintf("Product %q added.", opts.ProductName))

	return nil
}
