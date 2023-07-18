package server_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/server"
	"github.com/konstellation-io/kli/mocks"
)

type ListServersSuite struct {
	suite.Suite

	renderer *mocks.MockRenderer
	manager  *server.Handler
	tmpDir   string
}

func TestListServersSuite(t *testing.T) {
	suite.Run(t, new(ListServersSuite))
}

func (s *ListServersSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(logger)

	s.renderer = renderer
	s.manager = server.NewServerHandler(logger, renderer)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	err = os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), os.FileMode(0750))
	s.Require().NoError(err)

	s.tmpDir = tmpDir
}

func (s *ListServersSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ListServersSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ListServersSuite) TestListServers_ValidConfiguration() {
	storedConfig, err := createDefaultConfiguration()
	s.Require().NoError(err)

	s.renderer.EXPECT().RenderServers(storedConfig.Servers).Times(1)

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
	if err := os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), 0750); err != nil {
		return err
	}

	configBytes := []byte("invalid configuration")

	return os.WriteFile(viper.GetString(config.KaiConfigPath), configBytes, 0600)
}
