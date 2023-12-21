package render

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/karlseguin/jsonwriter"
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/krt/pkg/krt"
)

const (
	_no  = "No"
	_yes = "Yes"
)

type CliJSONRenderer struct {
	jsonWriter *jsonwriter.Writer
	ioWriter   io.Writer
}

func DefaultJSONWriter(w io.Writer) *jsonwriter.Writer {
	return jsonwriter.New(w)
}

func NewJSONRenderer(ioWriter io.Writer, jsonWriter *jsonwriter.Writer) *CliJSONRenderer {
	return &CliJSONRenderer{
		jsonWriter: jsonWriter,
		ioWriter:   ioWriter,
	}
}

func (r *CliJSONRenderer) RenderServers(servers []*configuration.Server) {
	r.jsonWriter.RootArray(func() {
		for _, s := range servers {
			defaultMark := ""
			if s.IsDefault {
				defaultMark = "*"
			}

			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Name", fmt.Sprintf("%s%s", s.Name, defaultMark))
				r.jsonWriter.KeyValue("URL", s.Host)
				r.jsonWriter.KeyValue("Authenticated", s.IsLoggedIn())
			})
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderWorkflows(workflows []krt.Workflow) {
	r.jsonWriter.RootArray(func() {
		for _, wf := range workflows {
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Name", wf.Name)
				r.jsonWriter.KeyValue("Type", string(wf.Type))
				r.jsonWriter.KeyValue("Configured properties", len(wf.Config))
				r.jsonWriter.KeyValue("Processes", len(wf.Processes))
			})
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderProcesses(processes []krt.Process) {
	r.jsonWriter.RootArray(func() {
		for _, pr := range processes { //nolint:gocritic
			gpu := _no

			if *pr.GPU {
				gpu = _yes
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

			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Name", pr.Name)
				r.jsonWriter.KeyValue("Type", string(pr.Type))
				r.jsonWriter.KeyValue("Image", pr.Image)
				r.jsonWriter.KeyValue("Replicas", *pr.Replicas)
				r.jsonWriter.KeyValue("GPU", gpu)
				r.jsonWriter.KeyValue("Subscriptions", strings.Join(pr.Subscriptions, "\n"))
				r.jsonWriter.KeyValue("Object store", obj)
				r.jsonWriter.KeyValue("CPU limits", fmt.Sprintf("Request: %s\nLimit: %s", pr.ResourceLimits.CPU.Request, pr.ResourceLimits.CPU.Limit))
				r.jsonWriter.KeyValue("Memory limits", fmt.Sprintf("Request: %s\nLimit: %s", pr.ResourceLimits.Memory.Request, pr.ResourceLimits.Memory.Limit)) //nolint:lll
				r.jsonWriter.KeyValue("Networking", net)
				r.jsonWriter.KeyValue("Configured properties", len(pr.Config))
				r.jsonWriter.KeyValue("Configured secrets", len(pr.Secrets))
			})
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderConfiguration(scope string, config map[string]string) {
	r.jsonWriter.RootObject(func() {
		for k, v := range config {
			r.jsonWriter.KeyValue(k, v)
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderRegisteredProcesses(registeredProcesses []*entity.RegisteredProcess) {
	r.jsonWriter.RootArray(func() {
		for _, pr := range registeredProcesses {
			r.jsonWriter.ArrayObject(func() {
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
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderProducts(products []kai.Product) {
	r.jsonWriter.RootArray(func() {
		for _, p := range products {
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("ID", p.ID)
				r.jsonWriter.KeyValue("Name", p.Name)
				r.jsonWriter.KeyValue("Description", p.Description)
			})
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderVersion(productID string, v *entity.Version) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Product", productID)
		r.jsonWriter.KeyValue("Tag", v.Tag)
		r.jsonWriter.KeyValue("Status", v.Status)
		r.jsonWriter.KeyValue("Creation Date", v.CreationDate.Format(time.RFC3339))
		if v.Error != "" {
			r.jsonWriter.KeyValue("Error", v.Error)
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderVersions(productID string, versions []*entity.Version) {
	r.jsonWriter.RootArray(func() {
		for _, v := range versions {
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Product", productID)
				r.jsonWriter.KeyValue("Tag", v.Tag)
				r.jsonWriter.KeyValue("Status", v.Status)
				r.jsonWriter.KeyValue("Creation Date", v.CreationDate.Format(time.RFC3339))
			})
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderTriggers(triggers []entity.TriggerEndpoint) {
	r.jsonWriter.RootArray(func() {
		for _, t := range triggers {
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Name", t.Trigger)
				r.jsonWriter.KeyValue("Status", "UP")
				r.jsonWriter.KeyValue("Link", t.URL)
			})
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderLogs(productID string, logs []entity.Log, _ entity.LogOutFormat, showAllLabels bool) {
	r.jsonWriter.RootArray(func() {
		for _, log := range logs {
			var fullLog string

			if !showAllLabels {
				fullLog = (log.FormatedLog)
			} else {
				fullLog = (fmt.Sprintf("%s - %s", log.FormatedLog, log.Labels))
			}

			r.jsonWriter.Value(fullLog)
		}
	})
	r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderCallout(_ *entity.Version) {
	// We do not need to render anything in JSON format.
}

func (r *CliJSONRenderer) RenderKliVersion(version, buildDate string) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Version", version)
		r.jsonWriter.KeyValue("Build Date", buildDate)
	})
	r.ioWriter.Write([]byte("\n"))
}
