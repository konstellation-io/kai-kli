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

func (kc *KaiConfiguration) UpdateServer(server Server) error {
	if !kc.CheckServerExists(server.Name) {
		return ErrServerNotFound
	}

	for i, serv := range kc.Servers {
		// Replace the server with the same name as the one to be updated.
		if serv.Name == server.Name {
			kc.Servers[i] = server
		}
	}

	// If the server to be updated is set as default, set the other servers to false.
	if server.IsDefault {
		if err := kc.SetDefaultServer(server.Name); err != nil {
			return err
		}
	}

	return nil
}

func (kc *KaiConfiguration) DeleteServer(server string) error {
	if !kc.CheckServerExists(server) {
		return ErrServerNotFound
	}

	var index int

	for i, serv := range kc.Servers {
		if serv.Name == server {
			index = i
			break
		}
	}

	redefineDefaultServer := kc.Servers[index].IsDefault

	kc.Servers = removeElementFromList(kc.Servers, index)

	if redefineDefaultServer && len(kc.Servers) > 0 {
		return kc.SetDefaultServer(kc.Servers[0].Name)
	}

	return nil
}

func (kc *KaiConfiguration) SetDefaultServer(serverName string) error {
	serverExists := kc.CheckServerExists(serverName)
	if !serverExists {
		return ErrServerNotFound
	}

	for i, s := range kc.Servers {
		s.IsDefault = serverName == s.Name
		kc.Servers[i] = s
	}

	return nil
}

func (kc *KaiConfiguration) CheckServerExists(serverName string) bool {
	for _, s := range kc.Servers {
		if serverName == s.Name {
			return true
		}
	}

	return false
}

func (kc *KaiConfiguration) GetServer(server string) (*Server, error) {
	for _, s := range kc.Servers {
		if server == s.Name {
			return &s, nil
		}
	}

	return nil, ErrServerNotFound
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

func removeElementFromList[T interface{}](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
