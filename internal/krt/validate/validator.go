package validate

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/konstellation-io/kli/internal/krt/entity"
	krterrors "github.com/konstellation-io/kli/internal/krt/errors"
	"github.com/konstellation-io/kli/internal/krt/fileutils"
	"github.com/konstellation-io/kli/internal/krt/fileutils/compression"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

type Validator struct {
	logger     logging.Interface
	compressor compression.Compressor
	renderer   render.Renderer
}

func NewValidator(logger logging.Interface, compressor compression.Compressor, renderer render.Renderer) *Validator {
	return &Validator{
		logger:     logger,
		compressor: compressor,
		renderer:   renderer,
	}
}

func (v *Validator) ValidateKrt(krtPath string) error {
	if path.Ext(krtPath) != ".krt" {
		return errors.New("input file is not a .krt file")
	}

	tmp, err := os.MkdirTemp("", "krt-validate-*")
	if err != nil {
		return fmt.Errorf("error creating temporal directory: %w", err)
	}

	defer os.RemoveAll(tmp)

	err = v.compressor.Extract(krtPath, tmp)
	if err != nil {
		return fmt.Errorf("error extracting krt file: %w", err)
	}

	krtYamlPath, _, err := fileutils.GetKrtYaml(tmp)
	if err != nil {
		return fmt.Errorf("error getting krt yaml path: %w", err)
	}

	yamlFile, err := entity.ParseFile(krtYamlPath)
	if err != nil {
		return fmt.Errorf("error parsing krt yaml: %w", err)
	}

	var validationErrors []*krterrors.ValidationError

	errs := yamlFile.Validate()
	if len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}

	errs = v.validateSrcPaths(yamlFile, tmp)
	if len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}

	protoError := v.validateEntrypointProto(yamlFile, tmp)
	if protoError != nil {
		validationErrors = append(validationErrors, protoError)
	}

	v.renderer.RenderValidationErrors(validationErrors)

	if validationErrors != nil {
		return errors.New("krt file is not valid")
	}

	return nil
}

func (v *Validator) validateSrcPaths(krtFile *entity.File, krtPath string) []*krterrors.ValidationError {
	var errs []*krterrors.ValidationError

	for _, workflow := range krtFile.Workflows {
		for _, node := range workflow.Nodes {
			nodeFile := path.Join(krtPath, node.Src)
			if !fileExists(nodeFile) {
				errs = append(errs, krterrors.NewSrcPathValidationError(string(workflow.Name), string(node.Name), node.Src))
			}
		}
	}

	return errs
}

func (v *Validator) validateEntrypointProto(krtYaml *entity.File, krtPath string) *krterrors.ValidationError {
	protoPath := path.Join(krtPath, krtYaml.Entrypoint.Proto)
	if !fileExists(protoPath) {
		return krterrors.NewValidationError("entrypoint.proto", "entrypoint proto file not found")
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
