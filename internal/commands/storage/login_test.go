//go:build unit

package storage_test

import (
	"fmt"
	"path"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/storage"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
	"github.com/konstellation-io/kli/pkg/osutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type StorageLoginSuite struct {
	suite.Suite

	handler    *storage.Handler
	testServer *configuration.Server
	logger     *mocks.MockLogger
}

func TestStorageLoginSuite(t *testing.T) {
	suite.Run(t, new(StorageLoginSuite))
}

func (s *StorageLoginSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	s.logger = logger

	s.handler = storage.NewHandler(logger)

	viper.Set(config.KaiConfigPath, path.Join("testdata", "kai.conf"))

	conf, err := configuration.NewKaiConfigService(logger).GetConfiguration()
	s.Require().NoError(err)

	s.testServer, err = conf.GetServerOrDefault("")
	s.Require().NoError(err)
}

func (s *StorageLoginSuite) TearDownSuite(_, _ string) {
	monkey.UnpatchAll()
}

func (s *StorageLoginSuite) TestStorageLogin() {
	patch := monkey.Patch(osutil.OpenBrowser, func(url string) error {
		s.Assert().Equal(s.testServer.StorageEndpoint, url)
		return nil
	})
	defer patch.Unpatch()

	err := s.handler.OpenConsole(s.testServer.Name)
	s.Require().NoError(err)
}

func (s *StorageLoginSuite) TestStorageLogin_CannotOpenBrowser_ExpectOk() {
	patch := monkey.Patch(osutil.OpenBrowser, func(url string) error {
		return fmt.Errorf("error opening browser")
	})
	defer patch.Unpatch()

	err := s.handler.OpenConsole(s.testServer.Name)
	s.Require().NoError(err)
}

func (s *StorageLoginSuite) TestStorageLogin_WrongServer_ExpectError() {
	err := s.handler.OpenConsole("wrong-name")
	s.Require().Error(err)
}
