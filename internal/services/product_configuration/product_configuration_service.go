package productconfiguration

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/konstellation-io/krt/pkg/parse"
	"github.com/spf13/viper"
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
	productConfigPath, err := GetProductConfigFilePath(product, productPath...)
	if err != nil {
		return nil, ErrProductConfigNotFound
	}

	configBytes, err := os.ReadFile(productConfigPath)
	if err != nil {
		return nil, ErrProductConfigNotFound
	}

	var productConfiguration krt.Krt

	err = yaml.Unmarshal(configBytes, &productConfiguration)
	if err != nil {
		return nil, ErrProductConfigNotFound
	}

	return &KaiProductConfiguration{
		Krt: &productConfiguration,
	}, nil
}

func (c *ProductConfigService) WriteConfiguration(newConfig *KaiProductConfiguration, product string, productPath ...string) error {
	if newConfig == nil {
		return ErrProductConfigNotFound
	}

	updatedConfig, err := parse.ParseKrtToYaml(newConfig.Krt)
	if err != nil {
		return err
	}

	productConfigPath, err := GetProductConfigFilePath(product, productPath...)
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

func GetProductConfigFilePath(product string, productPath ...string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current directory: %w", err)
	}

	productConfigPath := currentDir

	if len(productPath) != 0 {
		if path.IsAbs(productPath[0]) {
			productConfigPath = productPath[0]
		} else {
			productConfigPath = path.Join(currentDir, productPath[0])
		}
	}

	if strings.HasSuffix(productConfigPath, viper.GetString(config.KaiProductConfigFolder)) {
		return path.Join(productConfigPath, getProductKrt(product)), nil
	}

	return path.Join(productConfigPath, viper.GetString(config.KaiProductConfigFolder), getProductKrt(product)), nil
}

func getProductKrt(product string) string {
	product = strings.TrimSuffix(product, ".yaml")

	return fmt.Sprintf("%s.yaml",
		strings.ReplaceAll(strings.ReplaceAll(product, " ", "_"), ".", "_"))
}
