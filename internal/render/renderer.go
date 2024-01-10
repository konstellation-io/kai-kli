package render

import (
	"io"

	"github.com/konstellation-io/kli/cmd/config"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/spf13/viper"
)

type CliRenderer struct {
	Renderer
}

func NewDefaultCliRenderer(logger logging.Interface, w io.Writer) *CliRenderer {
	if viper.GetString(config.OutputFormatKey) == "json" {
		return &CliRenderer{NewJSONRenderer(logger, w, DefaultJSONWriter(w))}
	}

	return &CliRenderer{NewTextRenderer(logger, w, DefaultTableWriter(w))}
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
			fullLog = (log.FormatedLog)
		} else {
			fullLog = (fmt.Sprintf("%s - %s", log.FormatedLog, log.Labels))
		}

		_, err := fmt.Fprintln(file, fullLog)
		if err != nil {
			r.logger.Error("Error writing in logs.txt file: " + err.Error())
			return
		}
	}
}
