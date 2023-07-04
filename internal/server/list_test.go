package server_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/server"
	"github.com/konstellation-io/kli/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ListServersSuite struct {
	suite.Suite

	renderer *mocks.MockServerRenderer
	manager  *server.KaiConfigurator
	tmpDir   string
}

func TestListServersSuite(t *testing.T) {
	suite.Run(t, new(ListServersSuite))
}

func (s *ListServersSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockServerRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.renderer = renderer
	s.manager = server.NewKaiConfigurator(logger, renderer)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiPathKey, kaiConfigPath)

	s.tmpDir = tmpDir
}

func (s *ListServersSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ListServersSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiPathKey))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ListServersSuite) TestListServers_ValidConfiguration() {
	storedConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	s.renderer.EXPECT().RenderServers(storedConfig)

	err = s.manager.ListServers()
	s.Assert().NoError(err)
}

func (s *ListServersSuite) TestListServers_NotValidConfiguration() {
	err := createInvalidConfig()
	s.Require().NoError(err)

	err = s.manager.ListServers()
	unmarshalErr := &yaml.TypeError{}
	s.Assert().ErrorAs(err, &unmarshalErr)
}

func (s *ListServersSuite) TestListServers_NoConfigFound() {
	err := s.manager.ListServers()
	s.Assert().ErrorIs(err, os.ErrNotExist)
}

func createInvalidConfig() error {
	if err := os.Mkdir(path.Dir(viper.GetString(config.KaiPathKey)), 0750); err != nil {
		return err
	}

	configBytes := []byte("invalid configuration")

	return os.WriteFile(viper.GetString(config.KaiPathKey), configBytes, 0600)
}
