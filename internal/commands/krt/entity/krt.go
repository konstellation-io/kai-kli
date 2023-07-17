package entity

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ResourceName string

func NewResourceName(resourceName string) (ResourceName, error) {
	if IsValidResourceName(resourceName) {
		return ResourceName(resourceName), nil
	}

	return "", errors.New("invalid resource name")
}

// File contains data about a version.
type File struct {
	KrtVersion  string       `yaml:"krtVersion" validate:"required"`
	Version     ResourceName `yaml:"version" validate:"required,resource-name,lt=20"`
	Description string       `yaml:"description" validate:"required"`
	Entrypoint  *Entrypoint  `yaml:"entrypoint,omitempty" validate:"required"`
	Config      *Config      `yaml:"config,omitempty" validate:"required"`
	Workflows   []Workflow   `yaml:"workflows" validate:"required,dive,min=1"`
}

// Node contains data about a version's node.
type Node struct {
	Name          ResourceName `yaml:"name" validate:"required,resource-name,lt=20"`
	Image         string       `yaml:"image" validate:"required"`
	Src           string       `yaml:"src" validate:"required"`
	GPU           bool         `yaml:"gpu,omitempty" validate:"omitempty"`
	Replicas      int32        `yaml:"replicas,omitempty" validate:"omitempty"`
	Subscriptions []string     `yaml:"subscriptions" validate:"required"`
}

// Workflow contains data about a version's workflow.
type Workflow struct {
	Name       ResourceName `yaml:"name" validate:"required,resource-name,lt=20"`
	Entrypoint string       `yaml:"entrypoint" validate:"required"`
	Nodes      []Node       `yaml:"nodes" validate:"dive,min=1"`
	Exitpoint  string       `yaml:"exitpoint" validate:"required"`
}

// Entrypoint defines a KRT entrypoint Image and Proto file.
type Entrypoint struct {
	Image string `yaml:"image" validate:"required"`
	Proto string `yaml:"proto" validate:"required,endswith=.proto"`
}

// Config contains variables and file names.
type Config struct {
	Variables []string `yaml:"variables,omitempty" validate:"dive,env"`
	// TODO Check if this is used (or useful)
	Files []string `yaml:"files,omitempty" validate:"dive,env"`
}

// ParseKrt parse a Krt file from the given yaml bytes.
func ParseKrt(krtYaml []byte) (*File, error) {
	var file File

	err := yaml.Unmarshal(krtYaml, &file)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml file: %w", err)
	}

	return &file, nil
}

// ParseFile parse a Krt file from the given filename that must exists on the filesystem.
func ParseFile(yamlFile string) (*File, error) {
	var file File

	krtYmlFile, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}

	err = yaml.Unmarshal(krtYmlFile, &file)
	if err != nil {
		return nil, fmt.Errorf("error Unmarshal yaml file: %w", err)
	}

	return &file, nil
}
