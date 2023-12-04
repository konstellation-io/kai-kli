package version_test

import (
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	productName = "test-product"
	versionTag  = "v.1.0.1-test"
	comment     = "comment"
)

type VersionSuite struct {
	suite.Suite

	renderer         *mocks.MockRenderer
	logger           *mocks.MockLogger
	versionClient    *mocks.MockVersionClient
	handler          *version.Handler
	kaiConfiguration *configuration.KaiConfigService
	productConfig    *productconfiguration.ProductConfigService
	server           *configuration.Server
	tmpDir           string
	krtFilePath      string
}

func TestVersionSuite(t *testing.T) {
	suite.Run(t, new(VersionSuite))
}

func (s *VersionSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	s.renderer = mocks.NewMockRenderer(ctrl)
	s.versionClient = mocks.NewMockVersionClient(s.T())
	s.kaiConfiguration = configuration.NewKaiConfigService(s.logger)

	mocks.AddLoggerExpects(s.logger)

	tmpDir, err := os.MkdirTemp(s.T().TempDir(), "")
	s.Require().NoError(err)

	s.tmpDir = tmpDir

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	err = os.Mkdir(path.Dir(kaiConfigPath), os.FileMode(0750))
	s.Require().NoError(err)

	_, err = os.Create(kaiConfigPath)
	s.Require().NoError(err)

	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.server = &configuration.Server{Name: "test", Host: "http://test.com", IsDefault: true}

	kaiConf, err := s.kaiConfiguration.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(s.server)
	s.Require().NoError(err)

	err = s.kaiConfiguration.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	s.handler = version.NewHandler(s.logger, s.renderer, s.versionClient)
	s.krtFilePath = path.Join("testdata", productName+".yaml")

	s.productConfig = productconfiguration.NewProductConfigService(s.logger)
}

func (s *VersionSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}
