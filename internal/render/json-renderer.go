package render

import (
	"io"
	"time"

	"github.com/karlseguin/jsonwriter"
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/krt/pkg/krt"
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
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Name", s.Name)
				r.jsonWriter.KeyValue("URL", s.Host)
				r.jsonWriter.KeyValue("Authenticated", s.IsLoggedIn())
				r.jsonWriter.KeyValue("Default", s.IsDefault)
			})
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderWorkflows(workflows []krt.Workflow) {
	r.jsonWriter.RootArray(func() {
		for _, wf := range workflows {
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Name", wf.Name)
				r.jsonWriter.KeyValue("Type", string(wf.Type))
				r.jsonWriter.KeyValue("ConfiguredProperties", len(wf.Config))
				r.jsonWriter.KeyValue("Processes", len(wf.Processes))
			})
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderProcesses(processes []krt.Process) { //NOSONAR
	r.jsonWriter.RootArray(func() {
		for _, pr := range processes { //nolint:gocritic
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Name", pr.Name)
				r.jsonWriter.KeyValue("Type", string(pr.Type))
				r.jsonWriter.KeyValue("Image", pr.Image)
				r.jsonWriter.KeyValue("Replicas", *pr.Replicas)
				r.jsonWriter.KeyValue("GPU", *pr.GPU)
				r.jsonWriter.Array("Subscriptions", func() {
					for _, s := range pr.Subscriptions {
						r.jsonWriter.Value(s)
					}
				})
				r.jsonWriter.Object("ObjectStore", func() {
					if pr.ObjectStore == nil {
						r.jsonWriter.KeyValue("Name", "")
						r.jsonWriter.KeyValue("Scope", "")
					} else {
						r.jsonWriter.KeyValue("Name", pr.ObjectStore.Name)
						r.jsonWriter.KeyValue("Scope", pr.ObjectStore.Scope)
					}
				})
				r.jsonWriter.Object("CPULimits", func() {
					if pr.ResourceLimits.CPU == nil {
						r.jsonWriter.KeyValue("Request", "")
						r.jsonWriter.KeyValue("Limit", "")
					} else {
						r.jsonWriter.KeyValue("Request", pr.ResourceLimits.CPU.Request)
						r.jsonWriter.KeyValue("Limit", pr.ResourceLimits.CPU.Limit)
					}
				})
				r.jsonWriter.Object("MemoryLimits", func() {
					if pr.ResourceLimits.Memory == nil {
						r.jsonWriter.KeyValue("Request", "")
						r.jsonWriter.KeyValue("Limit", "")
					} else {
						r.jsonWriter.KeyValue("Request", pr.ResourceLimits.Memory.Request)
						r.jsonWriter.KeyValue("Limit", pr.ResourceLimits.Memory.Limit)
					}
				})
				r.jsonWriter.Object("Networking", func() {
					if pr.Networking == nil {
						r.jsonWriter.KeyValue("SourcePort", "")
						r.jsonWriter.KeyValue("ExposedPort", "")
						r.jsonWriter.KeyValue("Protocol", "")
					} else {
						r.jsonWriter.KeyValue("SourcePort", pr.Networking.TargetPort)
						r.jsonWriter.KeyValue("ExposedPort", pr.Networking.DestinationPort)
						r.jsonWriter.KeyValue("Protocol", string(pr.Networking.Protocol))
					}
				})
				r.jsonWriter.KeyValue("ConfiguredProperties", len(pr.Config))
				r.jsonWriter.KeyValue("ConfiguredSecrets", len(pr.Secrets))
			})
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderConfiguration(_ string, config map[string]string) {
	r.jsonWriter.RootObject(func() {
		for k, v := range config {
			r.jsonWriter.KeyValue(k, v)
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
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
				r.jsonWriter.KeyValue("UploadDate", pr.UploadDate.Format(time.RFC3339))
				r.jsonWriter.KeyValue("Owner", pr.Owner)
			})
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
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

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderVersion(productID string, v *entity.Version) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Product", productID)
		r.jsonWriter.KeyValue("Tag", v.Tag)
		r.jsonWriter.KeyValue("Status", v.Status)
		r.jsonWriter.KeyValue("CreationDate", v.CreationDate.Format(time.RFC3339))
		if v.Error != "" {
			r.jsonWriter.KeyValue("Error", v.Error)
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderVersions(productID string, versions []*entity.Version) {
	r.jsonWriter.RootArray(func() {
		for _, v := range versions {
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Product", productID)
				r.jsonWriter.KeyValue("Tag", v.Tag)
				r.jsonWriter.KeyValue("Status", v.Status)
				r.jsonWriter.KeyValue("CreationDate", v.CreationDate.Format(time.RFC3339))
			})
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
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

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderLogs(_ string, logs []entity.Log, _ entity.LogOutFormat, showAllLabels bool) {
	r.jsonWriter.RootArray(func() {
		for _, log := range logs {
			r.jsonWriter.ArrayObject(func() {
				r.jsonWriter.KeyValue("Message", log.FormatedLog)
				if showAllLabels {
					r.jsonWriter.Object("Labels", func() {
						for _, v := range log.Labels {
							r.jsonWriter.KeyValue(v.Key, v.Value)
						}
					})
				}
			})
		}
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderCallout(_ *entity.Version) {
	// We do not need to render anything in JSON format.
}

func (r *CliJSONRenderer) RenderKliVersion(version, buildDate string) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Version", version)
		r.jsonWriter.KeyValue("BuildDate", buildDate)
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}
