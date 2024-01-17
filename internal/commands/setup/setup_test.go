//go:build unit

package setup_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/commands/setup"
	"github.com/konstellation-io/kli/mocks"
)

type KliSetupTest struct {
	suite.Suite

	setup  setup.KaiSetup
	tmpDir string
}

func TestKliSetupSuite(t *testing.T) {
	suite.Run(t, new(KliSetupTest))
}

func (ks *KliSetupTest) SetupSuite() {
	ctrl := gomock.NewController(ks.T())
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	ks.setup = *setup.NewHandler(logger)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	ks.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	ks.tmpDir = tmpDir
}

func (ks *KliSetupTest) TearDownSuite(_, _ string) {
	err := os.RemoveAll(ks.tmpDir)
	ks.Require().NoError(err)
}

func (ks *KliSetupTest) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			ks.T().Fatalf("error cleaning tmp path: %ks", err)
		}
	}
}

func (ks *KliSetupTest) TestCreateConfiguration_ConfigFolderAndFileDoNotExist_ExpectOK() {
	// GIVEN
	ks.Require().NoFileExists(viper.GetString(config.KaiConfigPath))

	// WHEN
	err := ks.setup.CreateConfiguration()

	// THEN
	ks.Require().NoError(err)
	ks.Require().FileExists(viper.GetString(config.KaiConfigPath))
}

func (ks *KliSetupTest) TestCreateConfiguration_ConfigFolderAndFileAlreadyExist_ExpectOK() {
	// GIVEN
	err := ks.setup.CreateConfiguration()
	ks.Require().NoError(err)
	ks.Require().FileExists(viper.GetString(config.KaiConfigPath))
	file, err := os.ReadFile(viper.GetString(config.KaiConfigPath))
	ks.Require().NoError(err)

	// WHEN
	err = ks.setup.CreateConfiguration()

	// THEN
	ks.Require().NoError(err)
	ks.Require().FileExists(viper.GetString(config.KaiConfigPath))
	file2, err := os.ReadFile(viper.GetString(config.KaiConfigPath))
	ks.Require().NoError(err)
	ks.Require().Equal(file, file2)
}

func (ks *KliSetupTest) TestCheckConfigExists_ExpectTrue() {
	// GIVEN
	err := ks.setup.CreateConfiguration()
	ks.Require().NoError(err)
	ks.Require().FileExists(viper.GetString(config.KaiConfigPath))

	// WHEN
	fileExists := ks.setup.CheckConfigurationExists()

	// THEN
	ks.Require().True(fileExists)
}

func (ks *KliSetupTest) TestCheckConfigExists_ExpectFalse() {
	// GIVEN
	ks.Require().NoFileExists(viper.GetString(config.KaiConfigPath))

	// WHEN
	fileExists := ks.setup.CheckConfigurationExists()

	// THEN
	ks.Require().False(fileExists)
}
