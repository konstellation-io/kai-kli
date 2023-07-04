package server

import "errors"

var (
	ErrDuplicatedServerName = errors.New("duplicated server name")
	ErrDuplicatedServerURL  = errors.New("duplicated server URL")
	ErrServerNotFound       = errors.New("server not found")
)

type KaiConfiguration struct {
	DefaultServer string   `yaml:"default_server"`
	Servers       []Server `yaml:"servers"`
}

func (kc *KaiConfiguration) AddServer(server Server) error {
	if err := kc.checkServerDuplication(server); err != nil {
		return err
	}

	kc.Servers = append(kc.Servers, server)

	return nil
}

func (kc *KaiConfiguration) SetDefaultServer(serverName string) error {
	if !kc.checkServerExists(serverName) {
		return ErrServerNotFound
	}

	kc.DefaultServer = serverName

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
