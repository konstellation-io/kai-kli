package server

func (c *Handler) Logout(serverName string) error {
	return c.authentication.Logout(serverName)
}
