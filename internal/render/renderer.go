package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/gookit/color.v1"

	"github.com/konstellation-io/kli/api/kai/product"
	"github.com/konstellation-io/kli/api/kai/version"
	"github.com/konstellation-io/kli/internal/commands/krt/errors"
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
	}

	r.tableWriter.SetHeader([]string{
		"",
		"Name",
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
		"field",
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
	}

	productsNames := product.GetProductNames(products)
	r.logger.Success(fmt.Sprintf("%d products found: %s", len(products), strings.Join(productsNames, ", ")))
}

func (r *CliRenderer) RenderServers(servers []configuration.Server) {
	if len(servers) < 1 {
		r.logger.Info("No servers configured.")
		return
	}

	r.tableWriter.SetHeader([]string{"Server", "URL", "Auth"})

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
