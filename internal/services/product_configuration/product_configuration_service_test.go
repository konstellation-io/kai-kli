package productconfiguration_test

import (
	"errors"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ProductConfigurationServiceTest struct {
	suite.Suite

	productConfig productconfiguration.ProductConfigService
	productName   string
	tmpDir        string
}

func TestProductConfigurationServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductConfigurationServiceTest))
}

func (ch *ProductConfigurationServiceTest) SetupSuite() {
	ctrl := gomock.NewController(ch.T())
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	ch.productConfig = *productconfiguration.NewProductConfigService(logger)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	ch.Require().NoError(err)

	ch.tmpDir = tmpDir
	ch.productName = "test-product"
}

func (ch *ProductConfigurationServiceTest) TearDownSuite(_, _ string) {
	err := os.RemoveAll(ch.tmpDir)
	ch.Require().NoError(err)
}

func (ch *ProductConfigurationServiceTest) AfterTest(_, _ string) {
	if err := os.RemoveAll(ch.tmpDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			ch.T().Fatalf("error cleaning tmp path: %ch", err)
		}
	}
}

func (ch *ProductConfigurationServiceTest) TestReadProductConfig_ReadExistingConfiguration_ExpectOK() {
	// GIVEN
	defaultConfig, err := ch.createDefaultConfiguration()
	ch.Require().NoError(err)

	// WHEN
	currentConfig, err := ch.productConfig.GetConfiguration(ch.productName, ch.tmpDir)

	// THEN
	ch.Require().NoError(err)
	ch.NotNil(currentConfig)
	ch.Equal(defaultConfig.Version, currentConfig.Version)
	ch.Equal(defaultConfig.Description, currentConfig.Description)
	ch.True(reflect.DeepEqual(defaultConfig.Config, currentConfig.Config))

	for i, wf := range currentConfig.Workflows {
		ch.Equal(defaultConfig.Workflows[i].Name, wf.Name)
		ch.Equal(defaultConfig.Workflows[i].Type, wf.Type)
		ch.True(reflect.DeepEqual(defaultConfig.Workflows[i].Config, wf.Config))

		for j, pr := range wf.Processes { //nolint:gocritic
			defaultConfigProcess := defaultConfig.Workflows[i].Processes[j]
			ch.Equal(defaultConfigProcess.Name, pr.Name)
			ch.Equal(defaultConfigProcess.Type, pr.Type)
			ch.Equal(defaultConfigProcess.Image, pr.Image)
			ch.Equal(*defaultConfigProcess.Replicas, *pr.Replicas)
			ch.Equal(*defaultConfigProcess.GPU, *pr.GPU)
			ch.Equal(defaultConfigProcess.ObjectStore, pr.ObjectStore)
			ch.Equal(defaultConfigProcess.Subscriptions, pr.Subscriptions)
			ch.Empty(pr.Secrets)
			ch.True(reflect.DeepEqual(defaultConfigProcess.Config, pr.Config))
			ch.Equal(defaultConfigProcess.Networking, pr.Networking)
		}
	}
}

func (ch *ProductConfigurationServiceTest) TestReadProductConfig_ReadNonExistingConfiguration_ExpectError() {
	// GIVEN no config file
	// WHEN
	currentConfig, err := ch.productConfig.GetConfiguration(ch.productName)

	// THEN
	ch.Require().Error(err)
	ch.Nil(currentConfig)
}

func (ch *ProductConfigurationServiceTest) TestWriteProductConfig_WriteValidConfiguration_ExpectOK() {
	// GIVEN
	defaultConfig := ch.getDefaultConfiguration()

	// WHEN
	err := ch.productConfig.WriteConfiguration(&defaultConfig, ch.productName, ch.tmpDir)

	// THEN
	ch.Require().NoError(err)
	currentConfig, err := ch.productConfig.GetConfiguration(ch.productName, ch.tmpDir)
	ch.Require().NoError(err)
	ch.NotNil(currentConfig)
	ch.Equal(defaultConfig.Version, currentConfig.Version)
	ch.Equal(defaultConfig.Description, currentConfig.Description)
	ch.True(reflect.DeepEqual(defaultConfig.Config, currentConfig.Config))

	for i, wf := range currentConfig.Workflows {
		ch.Equal(defaultConfig.Workflows[i].Name, wf.Name)
		ch.Equal(defaultConfig.Workflows[i].Type, wf.Type)
		ch.True(reflect.DeepEqual(defaultConfig.Workflows[i].Config, wf.Config))

		for j, pr := range wf.Processes { //nolint:gocritic
			defaultConfigProcess := defaultConfig.Workflows[i].Processes[j]
			ch.Equal(defaultConfigProcess.Name, pr.Name)
			ch.Equal(defaultConfigProcess.Type, pr.Type)
			ch.Equal(defaultConfigProcess.Image, pr.Image)
			ch.Equal(*defaultConfigProcess.Replicas, *pr.Replicas)
			ch.Equal(*defaultConfigProcess.GPU, *pr.GPU)
			ch.Equal(defaultConfigProcess.ObjectStore, pr.ObjectStore)
			ch.Equal(defaultConfigProcess.Subscriptions, pr.Subscriptions)
			ch.Empty(pr.Secrets)
			ch.True(reflect.DeepEqual(defaultConfigProcess.Config, pr.Config))
			ch.Equal(defaultConfigProcess.Networking, pr.Networking)
		}
	}
}

func (ch *ProductConfigurationServiceTest) TestProductWriteConfig_WriteInvalidConfiguration_ExpectError() {
	// GIVEN a nil configuration
	// WHEN
	err := ch.productConfig.WriteConfiguration(nil, ch.productName, ch.tmpDir)

	// THEN
	ch.Require().Error(err)
	ch.Require().ErrorIs(err, productconfiguration.ErrProductConfigNotFound)
}

func (ch *ProductConfigurationServiceTest) getDefaultConfiguration() productconfiguration.KaiProductConfiguration {
	replicas := 1
	isGpu := false

	return productconfiguration.KaiProductConfiguration{
		Krt: &krt.Krt{
			Version:     "v0.0.1",
			Description: "KRT for Kai",
			Config:      map[string]string{"test1": "value1"},
			Workflows: []krt.Workflow{
				{
					Name:   "test",
					Type:   krt.WorkflowTypeData,
					Config: map[string]string{"test1": "value1"},
					Processes: []krt.Process{
						{
							Name:          "test",
							Type:          krt.ProcessTypeTrigger,
							Image:         "test/trigger-image",
							Replicas:      &replicas,
							GPU:           &isGpu,
							Config:        map[string]string{"test1": "value1"},
							ObjectStore:   nil,
							Secrets:       nil,
							Subscriptions: []string{"test1", "test2"},
							Networking: &krt.ProcessNetworking{
								TargetPort:      10000,
								DestinationPort: 15000,
								Protocol:        krt.NetworkingProtocolTCP,
							},
						},
					},
				},
			},
		},
	}
}

func (ch *ProductConfigurationServiceTest) createDefaultConfiguration() (*productconfiguration.KaiProductConfiguration, error) {
	defaultConfig := ch.getDefaultConfiguration()

	if err := ch.createConfigFile(defaultConfig); err != nil {
		return nil, err
	}

	return &defaultConfig, nil
}

func (ch *ProductConfigurationServiceTest) createConfigFile(cfg productconfiguration.KaiProductConfiguration) error {
	confPath, err := productconfiguration.GetProductConfigFilePath(ch.productName, ch.tmpDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(confPath), 0750); err != nil {
		return err
	}

	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(confPath, configBytes, 0600)
}
