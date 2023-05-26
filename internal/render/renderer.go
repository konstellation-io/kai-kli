package render

import (
	"fmt"
	"io"

	"github.com/konstellation-io/kli/api/kre/config"
	"github.com/konstellation-io/kli/api/kre/product"
	"github.com/konstellation-io/kli/api/kre/version"
	"github.com/konstellation-io/kli/internal/krt/errors"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/text"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/gookit/color.v1"
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

// RenderServerList add server information to the renderer and show it.
func (r *CliRenderer) RenderServerList(servers []config.ServerConfig, defaultServer string) {
	if len(servers) < 1 {
		r.logger.Info("No servers found.")
		return
	}

	r.tableWriter.SetHeader([]string{"Server", "URL"})

	for _, s := range servers {
		defaultMark := ""
		isDefault := text.Normalize(s.Name) == defaultServer

		if isDefault {
			defaultMark = "*"
		}

		r.tableWriter.Append([]string{
			fmt.Sprintf("%s%s", s.Name, defaultMark),
			s.URL,
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderVars(cfg *version.Config, showValues bool) {
	if len(cfg.Vars) == 0 {
		r.logger.Info("No config found.")
		return
	}

	r.renderVariables(cfg, showValues)

	r.printEmptyLine()

	if cfg.Completed {
		r.logger.Success("Version config complete")
	} else {
		r.logger.Warn("Version config incomplete")
	}
}

func (r *CliRenderer) printEmptyLine() {
	_, _ = fmt.Fprintln(r.ioWriter)
}

func (r *CliRenderer) renderVariables(cfg *version.Config, showValues bool) {
	h := []string{
		"",
		"TYPE",
		"KEY",
	}

	if showValues {
		h = append(h, "VALUE")
	}

	r.tableWriter.SetHeader(h)

	for i, v := range cfg.Vars {
		s := []string{
			fmt.Sprint(i + 1),
			string(v.Type),
			v.Key,
		}

		if showValues {
			s = append(s, v.Value)
		}

		r.tableWriter.Append(s)
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderVersions(versions version.List) {
	if len(versions) < 1 {
		r.logger.Info("No versions found.")
		return
	}

	r.tableWriter.SetHeader([]string{
		"",
		"ID",
		"Status",
	})

	// TODO export colors to a different file
	for _, rn := range versions {
		var statusIcon string

		switch rn.Status {
		case "STARTING", "STOPPING":
			statusIcon = color.RGB(83, 144, 254, false).Sprintf("⊙") //nolint:gomnd
		case "STARTED":
			statusIcon = color.RGB(109, 170, 72, false).Sprintf("⊙") //nolint:gomnd
		case "PUBLISHED":
			statusIcon = color.RGB(109, 170, 72, false).Sprintf("●") //nolint:gomnd
		case "STOPPED":
			statusIcon = color.RGB(163, 163, 163, false).Sprintf("○") //nolint:gomnd
		case "ERROR":
			statusIcon = color.RGB(229, 46, 61, false).Sprintf("∅") //nolint:gomnd
		default:
			statusIcon = color.Bold.Render("⊙")
		}

		r.tableWriter.Append([]string{
			statusIcon,
			rn.Name,
			rn.Status,
		})
	}

	r.tableWriter.Render()
}

func (r *CliRenderer) RenderValidationErrors(validationErrors []*errors.ValidationError) {
	if len(validationErrors) < 1 {
		r.logger.Success("File successfully validated!")
		return
	}

	r.tableWriter.SetHeader([]string{
		"",
		"Field",
		"Error",
	})

	for _, err := range validationErrors {
		r.tableWriter.Append([]string{
			fmt.Sprintf("[%s]", color.Danger.Render("✖")),
			err.Field,
			err.Message,
		})
	}

	r.tableWriter.Render()
	r.printEmptyLine()
}

func (r *CliRenderer) RenderProducts(products []product.Product) {
	if len(products) < 1 {
		r.logger.Info("No products found.")
		return
	}

	r.tableWriter.SetHeader([]string{
		"ID",
		"Name",
	})

	for _, productItem := range products {
		r.tableWriter.Append([]string{
			productItem.ID,
			productItem.Name,
		})
	}

	r.tableWriter.Render()
	r.printEmptyLine()
}
