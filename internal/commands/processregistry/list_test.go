//go:build unit

package processregistry_test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kli/cmd/config"
	processRegistryCMD "github.com/konstellation-io/kli/internal/commands/processregistry"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/kli/mocks"
)

type ListProcessSuite struct {
	suite.Suite

	logger             *mocks.MockLogger
	renderer           *mocks.MockRenderer
	processRegistryAPI *mocks.MockAPIClient
	manager            *processRegistryCMD.Handler
	configuration      *configuration.KaiConfigService
	tmpDir             string
}

func TestListProcessSuite(t *testing.T) {
	suite.Run(t, new(ListProcessSuite))
}

func (s *ListProcessSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.logger = mocks.NewMockLogger(ctrl)
	renderer := mocks.NewMockRenderer(ctrl)
	mocks.AddLoggerExpects(s.logger)

	viper.SetDefault(config.KaiProductConfigFolder, ".kai")

	s.renderer = renderer

	s.processRegistryAPI = mocks.NewMockAPIClient(ctrl)

	s.manager = processRegistryCMD.NewHandler(
		s.logger,
		renderer,
		s.processRegistryAPI,
	)

	tmpDir, err := os.MkdirTemp("", "TestAddServer_*")
	s.Require().NoError(err)

	kaiConfigPath := path.Join(tmpDir, ".kai", "kai.conf")
	viper.Set(config.KaiConfigPath, kaiConfigPath)

	s.tmpDir = tmpDir
	s.configuration = configuration.NewKaiConfigService(s.logger)
}

func (s *ListProcessSuite) TearDownSuite(_, _ string) {
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *ListProcessSuite) AfterTest(_, _ string) {
	if err := os.RemoveAll(path.Dir(viper.GetString(config.KaiConfigPath))); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.T().Fatalf("error cleaning tmp path: %s", err)
		}
	}
}

func (s *ListProcessSuite) BeforeTest(_, _ string) {
	err := os.Mkdir(path.Dir(viper.GetString(config.KaiConfigPath)), os.FileMode(0750))
	s.Require().NoError(err)

	_, err = os.Create(viper.GetString(config.KaiConfigPath))
	s.Require().NoError(err)

	kaiConf, err := s.configuration.GetConfiguration()
	s.Require().NoError(err)

	err = kaiConf.AddServer(&configuration.Server{Name: "test", Host: "http://test.com", IsDefault: true})
	s.Require().NoError(err)

	err = s.configuration.WriteConfiguration(kaiConf)
	s.Require().NoError(err)
}

func (s *ListProcessSuite) TestList() {
	// GIVEN
	productID := _productID
	serverName := _serverName
	processName := _processName
	version := _version
	processType := _processType

	listFilters := &processRegistryCMD.ListFilters{
		ProcessName: processName,
		Version:     version,
		ProcessType: krt.ProcessType(processType),
	}

	testRegisteredProcess := &entity.RegisteredProcess{
		ID:         "test-process-id",
		Name:       processName,
		Version:    version,
		Type:       processType,
		Image:      "test-image",
		UploadDate: time.Now().UTC(),
		Owner:      "test",
	}

	retrievedRegisteredProcesses := []*entity.RegisteredProcess{testRegisteredProcess}

	s.processRegistryAPI.EXPECT().List(gomock.Any(), productID, processName, version, processType).Return(
		retrievedRegisteredProcesses,
		nil,
	)
	s.renderer.EXPECT().RenderRegisteredProcesses(retrievedRegisteredProcesses)

	// WHEN
	err := s.manager.ListProcesses(&processRegistryCMD.ListProcessesOpts{
		ServerName: serverName,
		ProductID:  productID,
		Filters:    listFilters,
	})

	// THEN
	s.Require().NoError(err)
}

func (s *ListProcessSuite) TestList_APIError() {
	// GIVEN
	productID := _productID
	serverName := _serverName

	s.processRegistryAPI.EXPECT().List(gomock.Any(), productID, "", "", "").Return(
		nil,
		fmt.Errorf("mock error"),
	)

	// WHEN
	err := s.manager.ListProcesses(&processRegistryCMD.ListProcessesOpts{
		ServerName: serverName,
		ProductID:  productID,
	})

	// THEN
	s.Require().Error(err)
}

func (s *ListProcessSuite) TestList_InvalidTypeFilter_APIError() {
	// GIVEN
	productID := _productID
	serverName := _serverName
	processType := "some-invalid-type"

	filters := &processRegistryCMD.ListFilters{
		ProcessType: krt.ProcessType(processType),
	}

	// WHEN
	err := s.manager.ListProcesses(&processRegistryCMD.ListProcessesOpts{
		ServerName: serverName,
		ProductID:  productID,
		Filters:    filters,
	})

	// THEN
	s.Require().Error(err)
}

func (s *ListProcessSuite) TestList_InvalidServer_ExpectError() {
	// GIVEN
	productID := _productID
	serverName := "invalid-server"

	// WHEN
	err := s.manager.ListProcesses(&processRegistryCMD.ListProcessesOpts{
		ServerName: serverName,
		ProductID:  productID,
	})

	// THEN
	s.Require().Error(err)
}
