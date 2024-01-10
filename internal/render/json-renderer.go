package render

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/karlseguin/jsonwriter"
	"github.com/konstellation-io/kli/api/kai"
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/krt/pkg/krt"
)

type CliJSONRenderer struct {
	logger     logging.Interface
	jsonWriter *jsonwriter.Writer
	ioWriter   io.Writer
}

func DefaultJSONWriter(w io.Writer) *jsonwriter.Writer {
	return jsonwriter.New(w)
}

func NewJSONRenderer(logger logging.Interface, ioWriter io.Writer, jsonWriter *jsonwriter.Writer) *CliJSONRenderer {
	return &CliJSONRenderer{
		logger:     logger,
		jsonWriter: jsonWriter,
		ioWriter:   ioWriter,
	}
}

func (r *CliJSONRenderer) RenderServers(servers []*configuration.Server) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
			for _, s := range servers {
				r.jsonWriter.ArrayObject(func() {
					r.jsonWriter.KeyValue("Name", s.Name)
					r.jsonWriter.KeyValue("URL", s.Host)
					r.jsonWriter.KeyValue("Authenticated", s.IsLoggedIn())
					r.jsonWriter.KeyValue("Default", s.IsDefault)
				})
			}
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderWorkflows(workflows []krt.Workflow) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
			for _, wf := range workflows {
				r.jsonWriter.ArrayObject(func() {
					r.jsonWriter.KeyValue("Name", wf.Name)
					r.jsonWriter.KeyValue("Type", string(wf.Type))
					r.jsonWriter.KeyValue("ConfiguredProperties", len(wf.Config))
					r.jsonWriter.KeyValue("Processes", len(wf.Processes))
				})
			}
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderProcesses(processes []krt.Process) { // NOSONAR
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
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
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderConfiguration(_ string, config map[string]string) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Object("Data", func() {
			for k, v := range config {
				r.jsonWriter.KeyValue(k, v)
			}
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderRegisteredProcesses(registeredProcesses []*entity.RegisteredProcess) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
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
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderProducts(products []kai.Product) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
			for _, p := range products {
				r.jsonWriter.ArrayObject(func() {
					r.jsonWriter.KeyValue("ID", p.ID)
					r.jsonWriter.KeyValue("Name", p.Name)
					r.jsonWriter.KeyValue("Description", p.Description)
				})
			}
		})
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
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
			for _, v := range versions {
				r.jsonWriter.ArrayObject(func() {
					r.jsonWriter.KeyValue("Product", productID)
					r.jsonWriter.KeyValue("Tag", v.Tag)
					r.jsonWriter.KeyValue("Status", v.Status)
					r.jsonWriter.KeyValue("CreationDate", v.CreationDate.Format(time.RFC3339))
				})
			}
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderTriggers(triggers []entity.TriggerEndpoint) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
			for _, t := range triggers {
				r.jsonWriter.ArrayObject(func() {
					r.jsonWriter.KeyValue("Name", t.Trigger)
					r.jsonWriter.KeyValue("Status", "UP")
					r.jsonWriter.KeyValue("Link", t.URL)
				})
			}
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderLogs(productID string, logs []entity.Log, outFormat entity.LogOutFormat, showAllLabels bool) {
	if len(logs) == 0 {
		return
	}

	if outFormat == entity.OutFormatConsole {
		r.renderLogsConsole(logs, showAllLabels)
	} else {
		r.renderLogsFile(productID, logs, showAllLabels)
	}
}

func (r *CliJSONRenderer) renderLogsConsole(logs []entity.Log, showAllLabels bool) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Array("Data", func() {
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
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) renderLogsFile(productID string, logs []entity.Log, showAllLabels bool) { // NOSONAR
	fileName := fmt.Sprintf("%s-logs-%s.txt", productID, time.Now().Format("2006-01-02T15:04:05"))

	file, err := os.Create(fileName)
	if err != nil {
		r.logger.Error(ErrorCreatingLogsFile(err))
		return
	}

	defer file.Close()

	filewriter := jsonwriter.New(file)

	filewriter.RootObject(func() {
		filewriter.KeyValue("Status", "OK")
		filewriter.KeyValue("Message", "")
		filewriter.Array("Data", func() {
			for _, log := range logs {
				filewriter.ArrayObject(func() {
					filewriter.KeyValue("Message", log.FormatedLog)
					if showAllLabels {
						filewriter.Object("Labels", func() {
							for _, v := range log.Labels {
								filewriter.KeyValue(v.Key, v.Value)
							}
						})
					}
				})
			}
		})
	})
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

func (r *CliJSONRenderer) RenderProcessRegistered(process *entity.RegisteredProcess) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Object("Data", func() {
			r.jsonWriter.KeyValue("Name", process.Name)
			r.jsonWriter.KeyValue("Version", process.Version)
			r.jsonWriter.KeyValue("Status", process.Status)
			r.jsonWriter.KeyValue("Type", process.Type)
			r.jsonWriter.KeyValue("Image", process.Image)
			r.jsonWriter.KeyValue("UploadDate", process.UploadDate.Format(time.RFC3339))
			r.jsonWriter.KeyValue("Owner", process.Owner)
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderProductCreated(product string, server *configuration.Server, initLocal bool) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Object("Data", func() {
			r.jsonWriter.KeyValue("Product", product)
			r.jsonWriter.KeyValue("Server", server.Name)
			r.jsonWriter.KeyValue("InitLocal", initLocal)
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}

func (r *CliJSONRenderer) RenderProductBinded(product *kai.Product) {
	r.jsonWriter.RootObject(func() {
		r.jsonWriter.KeyValue("Status", "OK")
		r.jsonWriter.KeyValue("Message", "")
		r.jsonWriter.Object("Data", func() {
			r.jsonWriter.KeyValue("ID", product.ID)
			r.jsonWriter.KeyValue("Name", product.Name)
			r.jsonWriter.KeyValue("Description", product.Description)
		})
	})

	_, _ = r.ioWriter.Write([]byte("\n"))
}
