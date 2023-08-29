package render

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/olekukonko/tablewriter"

	"github.com/konstellation-io/kli/api/processregistry"
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
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetRowLine(true)
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

	r.tableWriter.SetHeader([]string{"Workflow Name", "Workflow type", "Configured properties", "Processes"})

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

	r.tableWriter.SetHeader([]string{"Process Name", "Process type", "Image", "Replicas", "GPU", "Subscriptions",
		"Object store", "CPU limits", "Memory limits", "Networking", "Configured properties", "Configured secrets"})

	for _, pr := range processes { //nolint:gocritic
		gpu := "No"

		if *pr.GPU {
			gpu = "Yes"
		}

		obj := "-"

		if pr.ObjectStore != nil {
			obj = fmt.Sprintf("ObjectStore: %s\nScope: %s", pr.ObjectStore.Name, pr.ObjectStore.Scope)
		}

		net := "-"

		if pr.Networking != nil {
			net = fmt.Sprintf("Source Port: %d\nExposed Port: %d\nProtocol: %s",
				pr.Networking.TargetPort, pr.Networking.DestinationPort, pr.Networking.Protocol)
		}

		r.tableWriter.Append([]string{
			pr.Name,
			string(pr.Type),
			pr.Image,
			strconv.Itoa(*pr.Replicas),
			gpu,
			strings.Join(pr.Subscriptions, "\n"),
			obj,
			fmt.Sprintf("Request: %s\nLimit: %s", pr.ResourceLimits.CPU.Request, pr.ResourceLimits.CPU.Limit),
			fmt.Sprintf("Request: %s\nLimit: %s", pr.ResourceLimits.Memory.Request, pr.ResourceLimits.Memory.Limit),
			net,
			fmt.Sprintf("%d", len(pr.Config)),
			fmt.Sprintf("%d", len(pr.Secrets)),
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderConfiguration(scope string, config map[string]string) {
	if len(config) < 1 {
		r.logger.Info("No configuration found.")
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

func (r *CliRenderer) RenderRegisteredProcesses(registeredProcesses []processregistry.RegisteredProcess) {
	if len(registeredProcesses) < 1 {
		r.logger.Info("No processes found.")
		return
	}

	r.tableWriter.SetHeader([]string{
		"Process Name", "Process Version", "Process type", "Image", "Upload Date", "Owner",
	})

	for _, pr := range registeredProcesses {
		r.tableWriter.Append([]string{
			pr.Name,
			pr.Version,
			pr.Type,
			pr.Image,
			pr.UploadDate.Format(time.RFC3339),
			pr.Owner,
		})
	}

	r.tableWriter.Render()
}
