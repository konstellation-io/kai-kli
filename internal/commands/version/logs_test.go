//go:build unit

package version_test

import (
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/logging"
)

func (s *VersionSuite) TestLogs() {
	const (
		productName = "test-product"
		versionTag  = "v.1.0.1-test"
	)

	filters := &entity.LogFilters{
		ProductID: productName,
	}

	logs := []entity.Log{
		{
			FormatedLog: "{'log': 'some formatted log'}",
			Labels: entity.LabelSet{
				{
					Key:   "product",
					Value: productName,
				},
				{
					Key:   "version",
					Value: versionTag,
				},
			},
		},
	}

	s.versionClient.EXPECT().Logs(s.server, filters).Return(logs, nil).Once()
	s.renderer.EXPECT().RenderLogs(productName, logs, (entity.LogOutFormat)(logging.OutputFormatText), true)

	err := s.handler.GetLogs(&version.GetLogsOpts{
		ServerName:    s.server.Name,
		OutFormat:     logging.OutputFormatText,
		LogFilters:    filters,
		ShowAllLabels: true,
	})
	s.Assert().NoError(err)
}

func (s *VersionSuite) TestLogs_WrongServer_ExpectError() {
	err := s.handler.GetLogs(&version.GetLogsOpts{
		ServerName:    "invalid-server",
		OutFormat:     logging.OutputFormatText,
		ShowAllLabels: true,
	})
	s.Assert().Error(err)
}
