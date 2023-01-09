package config

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/OpenPeeDeeP/xdg"

	"github.com/konstellation-io/kli/internal/logger"

	"github.com/konstellation-io/kli/pkg/errors"
	"github.com/konstellation-io/kli/text"
)

//nolint:gochecknoglobals
var (
	once sync.Once
	cfg  *Config
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
}

// ServerConfig contains data to represent a Konstellation server.
type ServerConfig struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	APIToken string `yaml:"token"`
}

// NewConfig will read the config.yml file from current user home.
func NewConfig() *Config {
	once.Do(createConfig)

	return cfg
}

// NewConfigTest same as NewConfig but used during tests to create a new file on each test.
func NewConfigTest() *Config {
	createConfig()
	return cfg
}

func createConfig() {
	log := logger.NewDefaultLogger()
	d := xdg.New("konstellation-io", "kli")

	cfg = &Config{
		Filename: path.Join(d.ConfigHome(), "config.yml"),
	}

	err := cfg.readFile()
	if err != nil {
		log.Error(fmt.Sprintf("error reading config: %s", err))
	}

	// Add default values
	cfg.DefaultRequestTimeout = DefaultRequestTimeout
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
