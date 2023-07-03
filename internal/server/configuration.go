package server

import "errors"

var (
	ErrDuplicatedServerName = errors.New("duplicated server name")
	ErrDuplicatedServerURL  = errors.New("duplicated server URL")
)

type KaiConfiguration struct {
	Servers []Server `yaml:"servers"`
}

func (kc *KaiConfiguration) AddServer(server Server) error {
	if err := kc.checkServerDuplication(server); err != nil {
		return err
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
