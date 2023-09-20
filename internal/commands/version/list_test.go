package version_test

import (
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	apiVersion "github.com/konstellation-io/kli/api/version"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ListVersionSuite struct {
	suite.Suite

	renderer         *mocks.MockRenderer
	logger           *mocks.MockLogger
	versionClient    *mocks.MockVersionClient
	handler          *version.Handler
	kaiConfiguration *configuration.KaiConfigService
	server           *configuration.Server
	tmpDir           string
}

func TestListVersionSuite(t *testing.T) {
	suite.Run(t, new(ListVersionSuite))
}

func (s *ListVersionSuite) SetupSuite() {
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

	s.server = &configuration.Server{Name: "test", URL: "http://test.com", IsDefault: true}

	kaiConf, err := s.kaiConfiguration.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(s.server)
	s.Require().NoError(err)

	err = s.kaiConfiguration.WriteConfiguration(kaiConf)
	s.Require().NoError(err)

	s.handler = version.NewHandler(s.logger, s.renderer, s.versionClient)
}

func (s *ListVersionSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ListVersionSuite) TestListVersion() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
	)

	versions := []*apiVersion.Version{
		{
			Tag:          versionTag,
			CreationDate: time.Now(),
			Status:       "CREATED",
		},
	}

	s.versionClient.EXPECT().List(s.server, productName).Return(versions, nil).Once()
	s.renderer.EXPECT().RenderVersions(productName, versions)

	err := s.handler.ListVersions(&version.ListVersionsOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
	})
	s.Assert().NoError(err)
}

func (s *ListVersionSuite) TestListVersion_ErrorIfClientFails() {
	productName := "test-product"
	expectedError := errors.New("client error")

	s.versionClient.EXPECT().List(s.server, productName).Return(nil, expectedError).Once()

	err := s.handler.ListVersions(&version.ListVersionsOpts{
		ServerName: s.server.Name,
		ProductID:  productName,
	})
	s.Assert().EqualError(err, expectedError.Error())
}
