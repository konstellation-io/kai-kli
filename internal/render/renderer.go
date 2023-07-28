package render

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/olekukonko/tablewriter"

	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

type CliRenderer struct {
	logger      logging.Interface
	tableWriter *tablewriter.Table
	ioWriter    io.Writer
}

func DefaultTableWriter(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetTablePadding(" ") // pad with spaces
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	return table
}

func NewCliRenderer(logger logging.Interface, ioWriter io.Writer, tableWriter *tablewriter.Table) *CliRenderer {
	return &CliRenderer{
		logger:      logger,
		tableWriter: tableWriter,
		ioWriter:    ioWriter,
	}
}

func NewDefaultCliRenderer(logger logging.Interface, writer io.Writer) *CliRenderer {
	return NewCliRenderer(logger, writer, DefaultTableWriter(writer))
}

func (r *CliRenderer) RenderServers(servers []configuration.Server) {
	if len(servers) < 1 {
		r.logger.Info("No servers configured.")
		return
	}

	r.tableWriter.SetHeader([]string{"Server", "URL", "Authenticated"})

	for _, s := range servers {
		defaultMark := ""

		if s.IsDefault {
			defaultMark = "*"
		}

		r.tableWriter.Append([]string{
			fmt.Sprintf("%s%s", s.Name, defaultMark),
			s.URL,
			fmt.Sprintf("%t", s.IsLoggedIn()),
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderWorkflows(workflows []krt.Workflow) {
	if len(workflows) < 1 {
		r.logger.Info("No workflows found.")
		return
	}

	r.tableWriter.SetHeader([]string{"Workflow Name", "Workflow type", "Configured properties", "Configured processes"})

	for _, wf := range workflows {
		r.tableWriter.Append([]string{
			wf.Name,
			string(wf.Type),
			fmt.Sprintf("%d", len(wf.Config)),
			fmt.Sprintf("%d", len(wf.Processes)),
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderProcesses(processes []krt.Process) {
	if len(processes) < 1 {
		r.logger.Info("No processes found.")
		return
	}

	r.tableWriter.SetHeader([]string{"Process Name", "Process type", "Image",
		"Replicas", "GPU", "Subscriptions", "Configured properties"})

	for _, pr := range processes {
		gpu := "No"

		if *pr.GPU {
			gpu = "Yes"
		}

		r.tableWriter.Append([]string{
			pr.Name,
			string(pr.Type),
			pr.Image,
			strconv.Itoa(*pr.Replicas),
			gpu,
			strings.Join(pr.Subscriptions, ","),
			fmt.Sprintf("%d", len(pr.Config)),
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderConfiguration(scope string, config map[string]string) {
	if len(config) < 1 {
		r.logger.Info("No processes found.")
		return
	}

	r.logger.Info(fmt.Sprintf("Rendering configuration for %q scope", scope))

	r.tableWriter.SetHeader([]string{"Key", "Value"})

	for key, value := range config {
		r.tableWriter.Append([]string{
			key,
			value,
		})
	}

	r.tableWriter.Render()
}
