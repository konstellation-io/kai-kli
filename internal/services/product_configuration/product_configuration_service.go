package productconfiguration

import (
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/logging"
)

var ErrProductConfigNotFound = errors.New("product config not found")

type ProductConfigService struct {
	logger logging.Interface
}

func NewProductConfigService(logger logging.Interface) *ProductConfigService {
	return &ProductConfigService{
		logger: logger,
	}
}

func (c *ProductConfigService) GetConfiguration(product string, productPath ...string) (*KaiProductConfiguration, error) {
	productConfigPath, err := config.GetProductConfigFilePath(product, productPath...)
	if err != nil {
		return nil, err
	}

	configBytes, err := os.ReadFile(productConfigPath)
	if err != nil {
		return nil, err
	}

	var productConfiguration KaiProductConfiguration

	err = yaml.Unmarshal(configBytes, &productConfiguration)
	if err != nil {
		return nil, err
	}

	return &productConfiguration, nil
}

func (c *ProductConfigService) WriteConfiguration(newConfig *KaiProductConfiguration, product string, productPath ...string) error {
	if newConfig == nil {
		return ErrProductConfigNotFound
	}

	updatedConfig, err := yaml.Marshal(newConfig)
	if err != nil {
		return err
	}

	productConfigPath, err := config.GetProductConfigFilePath(product, productPath...)
	if err != nil {
		return err
	}

	folderPermissions := 0750

	err = os.MkdirAll(path.Dir(productConfigPath), os.FileMode(folderPermissions))
	if err != nil {
		return err
	}

	filePermissions := 0600

	err = os.WriteFile(productConfigPath, updatedConfig, os.FileMode(filePermissions))
	if err != nil {
		return err
	}

	return nil
}
