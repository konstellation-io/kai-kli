package version

import (
	"errors"
	"time"

	"github.com/konstellation-io/kli/api"
	"github.com/konstellation-io/kli/internal/commands/version"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
	"github.com/spf13/cobra"
)

const (
	_fromFlag         = "from"
	_toFlag           = "to"
	_workflowNameFlag = "workflow-name"
	_processNameFlag  = "process-name"
	_requestIDFlag    = "request-id"
	_loglevelFlag     = "log-level"
	_loggerFlag       = "logger"
	_outFlag          = "out"
	_limitFlag        = "limit"
	_allLabelsFlag    = "all-labels"
)

// NewLogsCmd creates a new command to get the logs of a version of a product given its tag.
func NewLogsCmd(logger logging.Interface) *cobra.Command {
	nArgs := 2
	cmd := &cobra.Command{
		Use: "logs <product_id> <version_id> --from <from_timestamp> --to <to_timestamp> [opts...]",
		Annotations: map[string]string{
			"authenticated": "true",
		},
		Args:  cobra.ExactArgs(nArgs),
		Short: "Get version of a product given a tag",
		Example: `
    	$ kli product version logs <product_id> <version_tag> --from <from_timestamp> --to <to_timestamp> [opts...]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			productID := args[0]
			versionTag := args[1]

			lf, err := formLogFilters(logger, productID, versionTag, cmd)
			if err != nil {
				return err
			}

			out, err := cmd.Flags().GetString(_outFlag)
			if err != nil {
				out = "console"
			}
			if !entity.LogOutFormat(out).IsValid() {
				return errors.New("invalid logs output format")
			}

			allLabels, err := cmd.Flags().GetBool(_allLabelsFlag)
			if err != nil {
				allLabels = false
			}

			serverName, err := cmd.Flags().GetString("server")
			if err != nil {
				serverName = ""
			}

			r := render.NewDefaultCliRenderer(logger, cmd.OutOrStdout())

			err = version.NewHandler(logger, r, api.NewKaiClient().VersionClient()).
				GetLogs(&version.GetLogsOpts{
					ServerName:    serverName,
					LogFilters:    &lf,
					OutFormat:     entity.LogOutFormat(out),
					ShowAllLabels: allLabels,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(_fromFlag, "", "Required, timestamp with the begging of the query's time range.")
	cmd.Flags().String(_toFlag, "", "Required, timestamp with the end of the query's time range.")
	cmd.Flags().String(_workflowNameFlag, "", "Optional, filter results by workflow name.")
	cmd.Flags().String(_processNameFlag, "", "Optional, filter results by process name.")
	cmd.Flags().String(_requestIDFlag, "", "Optional, filter results by request ID.")
	cmd.Flags().String(_loglevelFlag, "", "Optional, filter results by log level.")
	cmd.Flags().String(_loggerFlag, "", "Optional, filter results by logger.")
	cmd.Flags().String(_outFlag, "console", "Optional, output format. One of: console|file. By default: console.")
	cmd.Flags().Int(_limitFlag, 0, "Optional, limit the number of results. By default: 100.")
	cmd.Flags().Bool(_allLabelsFlag, false, "Optional, show all labels attached to the logs. By default: false.")

	return cmd
}

func formLogFilters(logger logging.Interface, productID, versionTag string, cmd *cobra.Command) (entity.LogFilters, error) {
	from, err := cmd.Flags().GetString(_fromFlag)
	if err != nil {
		return entity.LogFilters{}, err
	}

	parsedFrom, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return entity.LogFilters{}, err
	}

	to, err := cmd.Flags().GetString(_toFlag)
	if err != nil {
		return entity.LogFilters{}, err
	}

	parsedTo, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return entity.LogFilters{}, err
	}

	workflowName, err := cmd.Flags().GetString(_workflowNameFlag)
	if err != nil {
		workflowName = ""
	}

	processName, err := cmd.Flags().GetString(_processNameFlag)
	if err != nil {
		processName = ""
	}

	requestID, err := cmd.Flags().GetString(_requestIDFlag)
	if err != nil {
		requestID = ""
	}

	loglevel, err := cmd.Flags().GetString(_loglevelFlag)
	if err != nil {
		loglevel = ""
	}

	loggerFilter, err := cmd.Flags().GetString(_loggerFlag)
	if err != nil {
		loggerFilter = ""
	}

	limit, err := cmd.Flags().GetInt(_limitFlag)
	if err != nil || limit == 0 {
		logger.Info("INFO: no limit detected, using default value (100)")

		limit = 100
	}

	lf := entity.LogFilters{
		ProductID:    productID,
		VersionTag:   versionTag,
		From:         parsedFrom,
		To:           parsedTo,
		WorkflowName: workflowName,
		ProcessName:  processName,
		RequestID:    requestID,
		Level:        loglevel,
		Logger:       loggerFilter,
		Limit:        limit,
	}

	return lf, nil
}
