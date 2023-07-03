package server

import (
	"errors"
	"regexp"
)

var (
	_validServerName     = regexp.MustCompile(`^[A-Za-z0-9\-_]+$`)
	ErrInvalidServerName = errors.New("invalid server name, only alphanumeric and hyphens are supported")
)

// Server contains data to represent a KAI server.
type Server struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Default bool   `yaml:"default"`
}

func (s Server) Validate() error {
	if !_validServerName.MatchString(s.Name) {
		return ErrInvalidServerName
	}

	return nil
}

type RemoteServerInfo struct {
	IsKaiServer bool
}
