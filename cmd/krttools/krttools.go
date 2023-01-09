package krttools

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/internal/builder"
	"github.com/konstellation-io/kli/internal/krt"
	"github.com/konstellation-io/kli/internal/logger"
)

var ErrYamlNotFound = errors.New("no yaml file found")

// KrtTooler interface with util functions.
type KrtTooler interface {
	Validate(krtFile string) error
	Build(src, target, version string) error
}

// KrtTools utils for validating and building krts.
type KrtTools struct {
	validator *krt.CustomValidator
	builder   builder.Builder
}

// NewKrtTools factory for KrtTools.
func NewKrtTools(loggerI logger.Interface) KrtTooler {
	validator := krt.NewCustomValidator(loggerI)

	return &KrtTools{
		validator: validator,
		builder:   builder.New(validator),
	}
}

// Validate validate util for krt command.
func (k *KrtTools) Validate(yamlPath string) error {
	krtFile, err := k.ParseYaml(yamlPath)
	if err != nil {
		return NewValidationError(err)
	}

	err = k.validator.Run(krtFile)
	if err != nil {
		return NewValidationError(err)
	}

	return nil
}

// Build utility for the krt command.
func (k *KrtTools) Build(src, target, version string) error {
	yamlPath, err := k.getKrtYamlPath(src)

	if err != nil {
		return NewBuildError(err)
	}

	krtFile, err := k.ParseYaml(yamlPath)
	if err != nil {
		return NewBuildError(err)
	}

	err = k.builder.UpdateVersion(krtFile, yamlPath, version)
	if err != nil {
		return NewBuildError(err)
	}

	err = k.validator.Run(krtFile)
	if err != nil {
		return NewBuildError(err)
	}

	err = k.builder.Build(src, target)

	return err
}

// getKrtYaml checks source directory for yaml files.
func (k *KrtTools) getKrtYamlPath(src string) (string, error) {
	ymls := []string{filepath.Join(src, "krt.yaml"), filepath.Join(src, "krt.yml")}

	for _, file := range ymls {
		if k.fileExists(file) {
			return file, nil
		}
	}

	return "", ErrYamlNotFound
}

func (k *KrtTools) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (k *KrtTools) ParseYaml(yamlPath string) (*krt.File, error) {
	r, err := os.Open(yamlPath)
	if err != nil {
		return nil, nil
	}

	var file *krt.File

	krtYmlFile, err := ioutil.ReadAll(r) //nolint:gocritic
	if err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}

	err = yaml.Unmarshal(krtYmlFile, &file)
	if err != nil {
		return nil, fmt.Errorf("error Unmarshal yaml file: %w", err)
	}

	return file, nil
}
