package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/konstellation-io/kli/internal/krt/entity"
	"github.com/konstellation-io/kli/internal/krt/fileutils"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

type Configurator struct {
	logger   logging.Interface
	renderer render.Renderer
}

func NewConfigurator(
	logger logging.Interface,
	renderer render.Renderer,
) *Configurator {
	return &Configurator{
		logger:   logger,
		renderer: renderer,
	}
}

func (c *Configurator) Validate(krtPath string) error {
	krtYaml, err := entity.ParseFile(krtPath)
	if err != nil {
		return err
	}

	validationErrors := krtYaml.Validate()
	if validationErrors != nil {
		c.renderer.RenderValidationErrors(validationErrors)
		return fmt.Errorf("%d error(s) found validating krt", len(validationErrors))
	}

	c.logger.Success("Krt file is valid.")

	return nil
}

func (c *Configurator) UpdateVersion(version, input, output string) error {
	inputFileInfo, err := os.Stat(input)
	if err != nil {
		return errors.New("invalid input folder")
	}

	if !inputFileInfo.IsDir() {
		return errors.New("input path should be a folder")
	}

	yamlPath, yamlFileName, err := fileutils.GetKrtYaml(input)
	if err != nil {
		return fmt.Errorf("krt.yaml not found in input folder: %w", err)
	}

	outputYaml := path.Join(output, yamlFileName)

	err = fileutils.UpdateVersion(yamlPath, outputYaml, version)
	if err != nil {
		return fmt.Errorf("error updating krt version: %w", err)
	}

	c.logger.Success("Krt.yaml version modified")

	return nil
}
