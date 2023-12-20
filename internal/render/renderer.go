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
	} else if viper.GetString("output") == "text" {
		return &CliRenderer{NewTextRenderer(logger, w, DefaultTableWriter(w))}
	}

	return nil
}
