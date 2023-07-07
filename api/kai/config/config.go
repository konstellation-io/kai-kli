package config

import (
	"fmt"
	"path"
	"time"

	"github.com/OpenPeeDeeP/xdg"

	"github.com/konstellation-io/kli/pkg/errors"
	"github.com/konstellation-io/kli/text"
)

const (
	// DefaultRequestTimeout time used to time out all requests to Konstellation APIs.
	DefaultRequestTimeout = 2 * time.Minute
)

// Config holds the configuration values for the application.
type Config struct {
	Filename              string
	DefaultRequestTimeout time.Duration  `yaml:"defaultRequestTimeout"`
	DefaultServer         string         `yaml:"defaultServer"`
	ServerList            []ServerConfig `yaml:"servers"`
	Debug                 bool
	BuildVersion          string
}

// ServerConfig contains data to represent a Konstellation server.
type ServerConfig struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	APIToken string `yaml:"token"`
}

// NewConfig will read the config.yml file from current user home.
func NewConfig(buildVersion string) (*Config, error) {
	cfg, err := loadConfigFromFile()
	if err != nil {
		return cfg, err
	}

	cfg.BuildVersion = buildVersion
	cfg.DefaultRequestTimeout = DefaultRequestTimeout

	return cfg, nil
}

func loadConfigFromFile() (*Config, error) {
	d := xdg.New("konstellation-io", "kli")

	cfg := &Config{
		Filename: path.Join(d.ConfigHome(), "config.yml"),
	}

	err := cfg.readFile()
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	return cfg, nil
}

// GetByServerName returns a ServerConfig for the given server name.
func (c *Config) GetByServerName(name string) *ServerConfig {
	n := text.Normalize(name)
	for _, s := range c.ServerList {
		if text.Normalize(s.Name) == n {
			return &s
		}
	}

	return nil
}

// AddServer adds a ServerConfig to the config file.
func (c *Config) AddServer(server ServerConfig) error {
	exists := c.GetByServerName(server.Name)
	if exists != nil {
		return errors.ErrServerAlreadyExists
	}

	c.ServerList = append(c.ServerList, server)

	return c.Save()
}

// SetDefaultServer marks a server name as default to be used when no server parameter is provided.
func (c *Config) SetDefaultServer(name string) error {
	n := text.Normalize(name)

	srv := c.GetByServerName(n)
	if srv == nil {
		return errors.ErrUnknownServerName
	}

	c.DefaultServer = n

	return c.Save()
}
