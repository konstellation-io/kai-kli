package version

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

func (c *Client) Logs(server *configuration.Server, filters entity.LogFilters) ([]entity.Log, error) {
	query := `
		query Logs($filters: LogFilters!) {
			logs(filters: $filters) {
				formatedLog
        labels {
            key
            value
        }
			}
		}
		`

	vars := map[string]interface{}{
		"filters": map[string]interface{}{
			"productID":    filters.ProductID,
			"versionTag":   filters.VersionTag,
			"from":         filters.From,
			"to":           filters.To,
			"limit":        filters.Limit,
			"workflowName": filters.WorkflowName,
			"processName":  filters.ProcessName,
			"requestID":    filters.RequestID,
			"level":        filters.Level,
			"logger":       filters.Logger,
		},
	}

	var respData struct {
		Logs []entity.Log `json:"logs"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.Logs, err
}
