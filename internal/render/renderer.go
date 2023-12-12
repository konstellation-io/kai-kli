package render

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/krt/pkg/krt"
	"github.com/olekukonko/tablewriter"
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

func (r *CliRenderer) RenderServers(servers []*configuration.Server) {
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
			s.Host,
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

func (r *CliRenderer) RenderRegisteredProcesses(registeredProcesses []*entity.RegisteredProcess) {
	if len(registeredProcesses) < 1 {
		r.logger.Info("No processes found.")
		return
	}

	r.tableWriter.SetHeader([]string{
		"Name", "Version", "Status", "Type", "Image", "Upload Date", "Owner", "isPublic",
	})

	for _, pr := range registeredProcesses {
		r.tableWriter.Append([]string{
			pr.Name,
			pr.Version,
			pr.Status,
			pr.Type,
			pr.Image,
			pr.UploadDate.Format(time.RFC3339),
			pr.Owner,
			strconv.FormatBool(pr.IsPublic),
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderProducts(products []kai.Product) {
	if len(products) == 0 {
		r.logger.Info("No products found.")
		return
	}

	r.tableWriter.SetHeader([]string{
		"ID", "Name", "Description",
	})

	for _, product := range products {
		r.tableWriter.Append([]string{
			product.ID,
			product.Name,
			product.Description,
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderVersions(productID string, versions []*entity.Version) {
	if len(versions) < 1 {
		r.logger.Info("No versions found.")
		return
	}

	r.tableWriter.SetHeader([]string{
		"Product", "Tag", "Status", "Creation Date",
	})

	for _, v := range versions {
		r.tableWriter.Append([]string{
			productID,
			v.Tag,
			v.Status,
			v.CreationDate.Format(time.RFC3339),
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderVersion(productID string, v *entity.Version) {
	if v.Error != "" {
		r.logger.Info(fmt.Sprintf("%s - %s status is: %s and has an error: %s.", productID, v.Tag, v.Status, v.Error))
		return
	}

	r.logger.Info(fmt.Sprintf("%s - %s status is: %s.", productID, v.Tag, v.Status))
}

func (r *CliRenderer) RenderTriggers(triggers []entity.TriggerEndpoint) {
	if len(triggers) == 0 {
		return
	}

	r.tableWriter.SetHeader([]string{
		"Name",
		"Status",
		"Link",
	})

	for _, trigger := range triggers {
		r.tableWriter.Append([]string{
			trigger.Trigger,
			"UP",
			trigger.URL,
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderLogs(productID string, logs []entity.Log, outFormat entity.LogOutFormat, showAllLabels bool) {
	if len(logs) == 0 {
		return
	}

	if outFormat == entity.OutFormatConsole {
		r.renderLogsConsole(logs, showAllLabels)
	} else {
		r.renderLogsFile(productID, logs, showAllLabels)
	}
}

func (r *CliRenderer) renderLogsConsole(logs []entity.Log, showAllLabels bool) {
	for _, log := range logs {
		if !showAllLabels {
			fmt.Println(log.FormatedLog)
		} else {
			fmt.Printf("%s - %s\n", log.FormatedLog, log.Labels)
		}
	}
}

func (r *CliRenderer) renderLogsFile(productID string, logs []entity.Log, showAllLabels bool) {
	fileName := fmt.Sprintf("%s-logs-%s.txt", productID, time.Now().Format("2006-01-02T15:04:05"))

	file, err := os.Create(fileName)
	if err != nil {
		r.logger.Error("Error creating logs.txt file: " + err.Error())
		return
	}
	defer file.Close()

	for _, log := range logs {
		var fullLog string

		if !showAllLabels {
			fullLog = log.FormatedLog
		} else {
			fullLog = fmt.Sprintf("%s - %s", log.FormatedLog, log.Labels)
		}

		_, err := fmt.Fprintln(file, fullLog)
		if err != nil {
			r.logger.Error("Error writing in logs.txt file: " + err.Error())
			return
		}
	}
}
