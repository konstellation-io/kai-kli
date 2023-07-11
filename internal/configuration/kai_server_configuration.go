package configuration

import (
	"errors"
)

var (
	ErrDuplicatedServerName = errors.New("duplicated server name")
	ErrDuplicatedServerURL  = errors.New("duplicated server URL")
	ErrServerNotFound       = errors.New("server not found")
)

func (kc *KaiConfiguration) AddServer(server Server) error {
	if err := kc.checkServerDuplication(server); err != nil {
		return err
	}

	kc.Servers = append(kc.Servers, server)

	// If there is just one server, set the first one as default.
	if len(kc.Servers) == 1 || server.IsDefault {
		if err := kc.SetDefaultServer(server.Name); err != nil {
			return err
		}
	}

	return nil
}

func (kc *KaiConfiguration) SetDefaultServer(serverName string) error {
	serverExists := kc.checkServerExists(serverName)
	if !serverExists {
		return ErrServerNotFound
	}

	for i, s := range kc.Servers {
		s.IsDefault = serverName == s.Name
		kc.Servers[i] = s
	}

	return nil
}

func (kc *KaiConfiguration) checkServerDuplication(server Server) error {
	for _, s := range kc.Servers {
		if s.Name == server.Name {
			return ErrDuplicatedServerName
		}

		if s.URL == server.URL {
			return ErrDuplicatedServerURL
		}
	}

	return nil
}

func (kc *KaiConfiguration) checkServerExists(serverName string) bool {
	for _, s := range kc.Servers {
		if serverName == s.Name {
			return true
		}
	}

	return false
}
