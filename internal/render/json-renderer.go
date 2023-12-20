package render

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/karlseguin/jsonwriter"
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/krt/pkg/krt"
)

type CliJsonRenderer struct {
	logger     logging.Interface
	jsonWriter *jsonwriter.Writer
	ioWriter   io.Writer
}

func DefaultJsonWriter(w io.Writer) *jsonwriter.Writer {
	return jsonwriter.New(w)
}

func NewJsonRenderer(logger logging.Interface, ioWriter io.Writer, jsonWriter *jsonwriter.Writer) *CliJsonRenderer {
	return &CliJsonRenderer{
		logger:     logger,
		jsonWriter: jsonWriter,
		ioWriter:   ioWriter,
	}
}

func (r *CliJsonRenderer) RenderServers(servers []*configuration.Server) {
	if len(servers) < 1 {
		r.logger.Info("No servers configured.")
		return
	}

	r.jsonWriter.RootObject(func() {
		for k, s := range servers {
			r.jsonWriter.Object(fmt.Sprintf("Server%d", k), func() {
				defaultMark := ""
				if s.IsDefault {
					defaultMark = "*"
				}

				r.jsonWriter.KeyValue("Name", fmt.Sprintf("%s%s", s.Name, defaultMark))
				r.jsonWriter.KeyValue("URL", s.Host)
				r.jsonWriter.KeyValue("Authenticated", s.IsLoggedIn())
			})
		}
	})
}

func (r *CliJsonRenderer) RenderWorkflows(workflows []krt.Workflow) {
	if len(workflows) < 1 {
		r.logger.Info("No workflows found.")
		return
	}

	r.jsonWriter.RootObject(func() {
		for k, wf := range workflows {
			r.jsonWriter.Object(fmt.Sprintf("Workflow%d", k), func() {
				r.jsonWriter.KeyValue("Name", wf.Name)
				r.jsonWriter.KeyValue("Type", string(wf.Type))
				r.jsonWriter.KeyValue("Configured properties", len(wf.Config))
				r.jsonWriter.KeyValue("Processes", len(wf.Processes))
			})
		}
	})
}

func (r *CliJsonRenderer) RenderProcesses(processes []krt.Process) {
	if len(processes) < 1 {
		r.logger.Info("No processes found.")
		return
	}

	r.jsonWriter.RootObject(func() {
		for k, pr := range processes {
			r.jsonWriter.Object(fmt.Sprintf("Process%d", k), func() {
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

				r.jsonWriter.KeyValue("Name", pr.Name)
				r.jsonWriter.KeyValue("Type", string(pr.Type))
				r.jsonWriter.KeyValue("Image", pr.Image)
				r.jsonWriter.KeyValue("Replicas", strconv.Itoa(*pr.Replicas))
				r.jsonWriter.KeyValue("GPU", gpu)
				r.jsonWriter.KeyValue("Subscriptions", strings.Join(pr.Subscriptions, "\n"))
				r.jsonWriter.KeyValue("Object store", obj)
				r.jsonWriter.KeyValue("CPU limits", fmt.Sprintf("Request: %s\nLimit: %s", pr.ResourceLimits.CPU.Request, pr.ResourceLimits.CPU.Limit))
				r.jsonWriter.KeyValue("Memory limits", fmt.Sprintf("Request: %s\nLimit: %s", pr.ResourceLimits.Memory.Request, pr.ResourceLimits.Memory.Limit))
				r.jsonWriter.KeyValue("Networking", net)
				r.jsonWriter.KeyValue("Configured properties", len(pr.Config))
				r.jsonWriter.KeyValue("Configured secrets", len(pr.Secrets))
			})
		}
	})
}

func (r *CliJsonRenderer) RenderConfiguration(scope string, config map[string]string) {
	if len(config) < 1 {
		r.logger.Info("No configuration found.")
		return
	}

	r.jsonWriter.RootObject(func() {
		r.jsonWriter.Object(scope, func() {
			for k, v := range config {
				r.jsonWriter.KeyValue(k, v)
			}
		})
	})
}

func (r *CliJsonRenderer) RenderRegisteredProcesses(registeredProcesses []*entity.RegisteredProcess) {
	if len(registeredProcesses) < 1 {
		r.logger.Info("No registered processes found.")
		return
	}

	r.jsonWriter.RootObject(func() {
		for k, pr := range registeredProcesses {
			r.jsonWriter.Object(fmt.Sprintf("RegisteredProcess%d", k), func() {
				r.jsonWriter.KeyValue("Name", pr.Name)
				r.jsonWriter.KeyValue("Version", pr.Version)
				r.jsonWriter.KeyValue("Status", pr.Status)
				r.jsonWriter.KeyValue("Type", pr.Type)
				r.jsonWriter.KeyValue("Image", pr.Image)
				r.jsonWriter.KeyValue("Upload Date", pr.UploadDate.Format(time.RFC3339))
				r.jsonWriter.KeyValue("Owner", pr.Owner)
			})
		}
	})
}

func (r *CliJsonRenderer) RenderProducts(products []kai.Product) {
	if len(products) < 1 {
		r.logger.Info("No products found.")
		return
	}

	r.jsonWriter.RootObject(func() {
		for k, p := range products {
			r.jsonWriter.Object(fmt.Sprintf("Product%d", k), func() {
				r.jsonWriter.KeyValue("ID", p.ID)
				r.jsonWriter.KeyValue("Name", p.Name)
				r.jsonWriter.KeyValue("Description", p.Description)
			})
		}
	})
}

func (r *CliJsonRenderer) RenderVersion(productID string, v *entity.Version) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.Object("Version", func() {
			if v.Error != "" {
				r.jsonWriter.KeyValue("Message", fmt.Sprintf("%s - %s status is: %s and has an error: %s.", productID, v.Tag, v.Status, v.Error))
			} else {
				r.jsonWriter.KeyValue("Message", fmt.Sprintf("%s - %s status is: %s.", productID, v.Tag, v.Status))
			}
		})
	})
}

func (r *CliJsonRenderer) RenderVersions(productID string, versions []*entity.Version) {
	if len(versions) < 1 {
		r.logger.Info("No versions found.")
		return
	}

	r.jsonWriter.RootObject(func() {
		for k, v := range versions {
			r.jsonWriter.Object(fmt.Sprintf("Version%d", k), func() {
				r.jsonWriter.KeyValue("Product", productID)
				r.jsonWriter.KeyValue("Tag", v.Tag)
				r.jsonWriter.KeyValue("Status", v.Status)
				r.jsonWriter.KeyValue("Creation Date", v.CreationDate.Format(time.RFC3339))
			})
		}
	})
}

func (r *CliJsonRenderer) RenderTriggers(triggers []entity.TriggerEndpoint) {
	if len(triggers) < 1 {
		r.logger.Info("No triggers found.")
		return
	}

	r.jsonWriter.RootObject(func() {
		for k, t := range triggers {
			r.jsonWriter.Object(fmt.Sprintf("Trigger%d", k), func() {
				r.jsonWriter.KeyValue("Name", t.Trigger)
				r.jsonWriter.KeyValue("Status", "UP")
				r.jsonWriter.KeyValue("Link", t.URL)
			})
		}
	})
}

func (r *CliJsonRenderer) RenderLogs(productID string, logs []entity.Log, outFormat entity.LogOutFormat, showAllLabels bool) {
	logName := fmt.Sprintf("%s-logs-%s", productID, time.Now().Format("2006-01-02T15:04:05"))

	r.jsonWriter.RootObject(func() {
		for _, log := range logs {
			var fullLog string

			if !showAllLabels {
				fullLog = (log.FormatedLog)
			} else {
				fullLog = (fmt.Sprintf("%s - %s", log.FormatedLog, log.Labels))
			}

			r.jsonWriter.Object(logName, func() {
				r.jsonWriter.KeyValue("Log", fullLog)
			})
		}
	})
}
