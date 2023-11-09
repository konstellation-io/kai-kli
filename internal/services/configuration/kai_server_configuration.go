package configuration

import (
	"errors"
)

var (
	ErrDuplicatedServerName = errors.New("duplicated server name")
	ErrDuplicatedServerHost = errors.New("duplicated server host")
	ErrServerNotFound       = errors.New("server not found")
	ErrNoServersFound       = errors.New("there are no servers configured")
)

func (kc *KaiConfiguration) AddServer(server *Server) error {
	if err := kc.checkServerDuplication(server); err != nil {
		return err
	}

	kc.Servers = append(kc.Servers, *server)

	// If there is just one server, set the first one as default.
	if len(kc.Servers) == 1 || server.IsDefault {
		if err := kc.SetDefaultServer(server.Name); err != nil {
			return err
		}
	}

	return nil
}

func (kc *KaiConfiguration) UpdateServer(server *Server) error {
	index, _, err := kc.findServer(server.Name)
	if err != nil {
		return ErrServerNotFound
	}

	kc.Servers[index] = *server

	// If the server to be updated is set as default, set the other servers to false.
	if server.IsDefault || len(kc.Servers) == 1 {
		if err := kc.SetDefaultServer(server.Name); err != nil {
			return err
		}
	}

	return nil
}

func (kc *KaiConfiguration) DeleteServer(serverName string) error {
	index, server, err := kc.findServer(serverName)
	if err != nil {
		return ErrServerNotFound
	}

	kc.Servers = removeElementFromList(kc.Servers, index)

	if server.IsDefault && len(kc.Servers) > 0 {
		return kc.SetDefaultServer(kc.Servers[0].Name)
	}

	return nil
}

func (kc *KaiConfiguration) SetDefaultServer(serverName string) error {
	if !kc.CheckServerExists(serverName) {
		return ErrServerNotFound
	}

	for i, s := range kc.Servers {
		s.IsDefault = serverName == s.Name
		kc.Servers[i] = s
	}

	return nil
}

func (kc *KaiConfiguration) GetDefaultServer() (*Server, error) {
	if len(kc.Servers) == 0 {
		return nil, ErrNoServersFound
	}

	for _, s := range kc.Servers {
		if s.IsDefault {
			return &s, nil
		}
	}

	return nil, ErrServerNotFound
}

func (kc *KaiConfiguration) CheckServerExists(serverName string) bool {
	_, _, err := kc.findServer(serverName)
	return err == nil
}

func (kc *KaiConfiguration) GetServer(serverName string) (*Server, error) {
	_, server, err := kc.findServer(serverName)
	if err != nil {
		return nil, ErrServerNotFound
	}

	return server, nil
}

func (kc *KaiConfiguration) GetServerOrDefault(serverName string) (*Server, error) {
	if serverName == "" {
		return kc.GetDefaultServer()
	}

	_, server, err := kc.findServer(serverName)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (kc *KaiConfiguration) findServer(server string) (int, *Server, error) {
	for i, s := range kc.Servers {
		if server == s.Name {
			return i, &s, nil
		}
	}

	return 0, nil, ErrServerNotFound
}

func (kc *KaiConfiguration) checkServerDuplication(server *Server) error {
	for _, s := range kc.Servers {
		if s.Name == server.Name {
			return ErrDuplicatedServerName
		}

		if s.Host == server.Host {
			return ErrDuplicatedServerHost
		}
	}

	return nil
}

func removeElementFromList[T interface{}](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
