package krt

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// File contains data about a version.
type File struct {
	KrtVersion  string     `yaml:"krtVersion" validate:"required"`
	Version     string     `yaml:"version" validate:"required,resource-name,lt=20"`
	Description string     `yaml:"description" validate:"required"`
	Entrypoint  Entrypoint `yaml:"entrypoint" validate:"required"`
	Config      Config     `yaml:"config" validate:"required"`
	Workflows   []Workflow `yaml:"workflows" validate:"required,dive,min=1"`
}

// Node contains data about a version's node.
type Node struct {
	Name          string   `yaml:"name" validate:"required,resource-name,lt=20"`
	Image         string   `yaml:"image" validate:"required"`
	Src           string   `yaml:"src" validate:"required"`
	GPU           bool     `yaml:"gpu,omitempty" validate:"omitempty"`
	Subscriptions []string `yaml:"subscriptions" validate:"required"`
}

// Workflow contains data about a version's workflow.
type Workflow struct {
	Name       string `yaml:"name" validate:"required,resource-name,lt=20"`
	Entrypoint string `yaml:"entrypoint" validate:"required"`
	Nodes      []Node `yaml:"nodes" validate:"dive,min=1"`
	Exitpoint  string `yaml:"exitpoint" validate:"required"`
}

// Entrypoint defines a KRT entrypoint Image and Proto file.
type Entrypoint struct {
	Image string `yaml:"image" validate:"required"`
	Proto string `yaml:"proto" validate:"required,endswith=.proto"`
}

// Config contains variables and file names.
type Config struct {
	Variables []string `yaml:"variables,omitempty" validate:"dive,env"`
	Files     []string `yaml:"files,omitempty" validate:"dive,env"`
}

// ParseFile parse a Krt file from the given filename that must exists on the filesystem.
func ParseFile(yamlFile string) (*File, error) {
	reader, err := os.Open(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", yamlFile, err)
	}

	var file File

	krtYmlFile, err := ioutil.ReadAll(reader) //nolint:gocritic
	if err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}

	err = yaml.Unmarshal(krtYmlFile, &file)
	if err != nil {
		return nil, fmt.Errorf("error Unmarshal yaml file: %w", err)
	}

	return &file, nil
}
