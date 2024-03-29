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

type CliTextRenderer struct {
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

func NewTextRenderer(logger logging.Interface, ioWriter io.Writer, tableWriter *tablewriter.Table) *CliTextRenderer {
	return &CliTextRenderer{
		logger:      logger,
		tableWriter: tableWriter,
		ioWriter:    ioWriter,
	}
}

func (r *CliTextRenderer) RenderServers(servers []*configuration.Server) {
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

func (r *CliTextRenderer) RenderWorkflows(workflows []krt.Workflow) {
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

func (r *CliTextRenderer) RenderProcesses(processes []krt.Process) {
	if len(processes) < 1 {
		r.logger.Info("No processes found.")
		return
	}

	r.tableWriter.SetHeader([]string{"Process Name", "Process type", "Image", "Replicas", "GPU", "Subscriptions",
		"Object store", "CPU limits", "Memory limits", "Networking", "Configured properties", "Configured secrets"})

	for _, pr := range processes { //nolint:gocritic
		obj := "-"

		if pr.ObjectStore != nil {
			obj = fmt.Sprintf("ObjectStore: %s\nScope: %s", pr.ObjectStore.Name, pr.ObjectStore.Scope)
		}

		net := "-"

		if pr.Networking != nil {
			net = fmt.Sprintf("Source Port: %d\nExposed Port: %d\nProtocol: %s",
				pr.Networking.TargetPort, pr.Networking.DestinationPort, pr.Networking.Protocol)
		}

		replicas := 1
		if pr.Replicas != nil {
			replicas = *pr.Replicas
		}

		gpu := false
		if pr.GPU != nil {
			gpu = *pr.GPU
		}

		r.tableWriter.Append([]string{
			pr.Name,
			string(pr.Type),
			pr.Image,
			strconv.Itoa(replicas),
			boolToText(gpu),
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

func (r *CliTextRenderer) RenderConfiguration(scope string, config map[string]string) {
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

func (r *CliTextRenderer) RenderRegisteredProcesses(registeredProcesses []*entity.RegisteredProcess) {
	if len(registeredProcesses) < 1 {
		r.logger.Info("No registered processes found.")
		return
	}

	r.tableWriter.SetHeader([]string{
		"Name", "Version", "Status", "Type", "Image", "Upload Date", "Owner", "Is Public",
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
			boolToText(pr.IsPublic),
		})
	}

	r.tableWriter.Render()
}

func (r *CliTextRenderer) RenderProducts(products []kai.Product) {
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

func (r *CliTextRenderer) RenderVersion(productID string, v *entity.Version) {
	if v.Error != "" {
		r.logger.Info(fmt.Sprintf("%s - %s status is: %s and has an error: %s.", productID, v.Tag, v.Status, v.Error))
		return
	}

	r.logger.Info(fmt.Sprintf("%s - %s status is: %s.", productID, v.Tag, v.Status))

	if v.Status == entity.VersionStatusPublished {
		r.logger.Info("Published triggers:")
		r.renderTriggers(v.PublishedTriggers)
	}
}

func (r *CliTextRenderer) RenderVersions(productID string, versions []*entity.Version) {
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

func (r *CliTextRenderer) RenderLogs(productID string, logs []entity.Log, outFormat entity.LogOutFormat, showAllLabels bool) {
	if len(logs) == 0 {
		return
	}

	if outFormat == entity.OutFormatConsole {
		r.renderLogsConsole(logs, showAllLabels)
	} else {
		r.renderLogsFile(productID, logs, showAllLabels)
	}
}

func (r *CliTextRenderer) renderLogsConsole(logs []entity.Log, showAllLabels bool) {
	for _, log := range logs {
		if !showAllLabels {
			fmt.Println(log.FormatedLog)
		} else {
			fmt.Printf("%s - %s\n", log.FormatedLog, log.Labels)
		}
	}
}

func (r *CliTextRenderer) renderLogsFile(productID string, logs []entity.Log, showAllLabels bool) {
	fileName := fmt.Sprintf("%s-logs-%s.txt", productID, time.Now().Format("2006-01-02T15:04:05"))

	file, err := os.Create(fileName)
	if err != nil {
		r.logger.Error(ErrorCreatingLogsFile(err))
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

func (r *CliTextRenderer) RenderCallout(v *entity.Version) {
	fmt.Println(v)
	r.logger.Warn("WARNING: Version is in CRITICAL state, start under your own discretion.")
}

func (r *CliTextRenderer) RenderKliVersion(version, buildDate string) {
	version = strings.TrimPrefix(version, "v")

	if buildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, buildDate)
	}

	fmt.Printf("kli version %s\n", version)
}

func (r *CliTextRenderer) RenderProcessRegistered(process *entity.RegisteredProcess) {
	r.logger.Success(fmt.Sprintf("Creating process with id %q, image will be uploaded to %q", process.ID, process.Image))
}

func (r *CliTextRenderer) RenderProcessDeleted(process string) {
	r.logger.Success(fmt.Sprintf("Process with id %q successfully deleted!", process))
}

func (r *CliTextRenderer) RenderProductCreated(_ string, _ *configuration.Server, _ bool) {
	// We do not need to render anything in text format.
}

func (r *CliTextRenderer) RenderProductBinded(productID string) {
	r.logger.Success(fmt.Sprintf("Product %q successfully bound!", productID))
}

func (r *CliTextRenderer) RenderLogin(serverName string) {
	r.logger.Success(fmt.Sprintf("Logged into %q.", serverName))
}

func (r *CliTextRenderer) RenderLogout(serverName string) {
	r.logger.Success(fmt.Sprintf("Logged out from %q.", serverName))
}

func (r *CliTextRenderer) RenderPushVersion(product, versionTag string) {
	r.logger.Success(fmt.Sprintf("Version with tag %q of product %q successfully created!", versionTag, product))
}

func (r *CliTextRenderer) RenderStartVersion(product, versionTag string) {
	r.logger.Success(fmt.Sprintf("Version with tag %q of product %q is starting!", versionTag, product))
}

func (r *CliTextRenderer) RenderStopVersion(product, versionTag string) {
	r.logger.Success(fmt.Sprintf("Version with tag %q of product %q is stopping!", versionTag, product))
}

func (r *CliTextRenderer) RenderPublishVersion(product, versionTag string, triggers []entity.TriggerEndpoint) {
	r.logger.Success(fmt.Sprintf("Version with tag %q of product %q correctly published!", versionTag, product))
	r.renderTriggers(triggers)
}

func (r *CliTextRenderer) renderTriggers(triggers []entity.TriggerEndpoint) {
	if len(triggers) == 0 {
		r.logger.Info("No triggers found.")
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

func (r *CliTextRenderer) RenderUnpublishVersion(product, versionTag string) {
	r.logger.Success(fmt.Sprintf("Version with tag %q of product %q correctly unpublished!", versionTag, product))
}

func boolToText(b bool) string {
	if b {
		return "Yes"
	}

	return "No"
}

func (r *CliTextRenderer) RenderAddUserToProduct(product, userEmail string) {
	r.logger.Success(fmt.Sprintf("User with email %q added to product %q!", userEmail, product))
}

func (r *CliTextRenderer) RenderRemoveUserFromProduct(product, userEmail string) {
	r.logger.Success(fmt.Sprintf("User with email %q removed from product %q!", userEmail, product))
}

func (r *CliTextRenderer) RenderAddMaintainerToProduct(product, userEmail string) {
	r.logger.Success(fmt.Sprintf("Maintainer with email %q added to product %q!", userEmail, product))
}

func (r *CliTextRenderer) RenderRemoveMaintainerFromProduct(product, userEmail string) {
	r.logger.Success(fmt.Sprintf("Maintainer with email %q removed from product %q!", userEmail, product))
}
