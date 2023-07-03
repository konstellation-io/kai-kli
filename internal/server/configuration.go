package server

import "errors"

var (
	ErrDuplicatedServerName = errors.New("duplicated server name")
	ErrDuplicatedServerURL  = errors.New("duplicated server URL")
)

type KaiConfiguration struct {
	DefaultServer string   `yaml:"default_server"`
	Servers       []Server `yaml:"servers"`
}

func (kc *KaiConfiguration) AddServer(server Server) error {
	if err := kc.checkServerDuplication(server); err != nil {
		return err
	}

	if server.Default {
		kc.DefaultServer = server.Name
	}

	kc.Servers = append(kc.Servers, server)

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
