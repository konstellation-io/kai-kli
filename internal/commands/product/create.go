package product

import (
	"errors"
	"os"

	"github.com/konstellation-io/krt/pkg/krt"

	"github.com/konstellation-io/kli/cmd/config"
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

var ErrProductConfigurationAlreadyExists = errors.New("product configuration already exists")

func (c *Handler) CreateProduct( /*server configuration.Server, */ productName, version, description string,

/*initLocal bool, localPath string*/) error {
	// TODO create product in KAI server

	productConfigPath, err := config.GetProductConfigFilePath(productName)
	if err != nil {
		return err
	}

	_, err = os.Stat(productConfigPath)
	if !os.IsNotExist(err) {
		return ErrProductConfigurationAlreadyExists
	}

	err = c.configService.WriteConfiguration(
		&productconfiguration.KaiProductConfiguration{
			Krt: &krt.Krt{
				Version:     productName,
				Description: description,
				Config:      map[string]string{"test": "value"},
				Workflows:   []krt.Workflow{},
			},
		}, productName)
	if err != nil {
		return err
	}

	return nil
}
