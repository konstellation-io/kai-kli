package productconfiguration

import (
	"errors"
	"golang.org/x/mod/semver"
)

var (
	ErrInvalidVersion = errors.New("invalid version format")
)

func (c *KaiProductConfiguration) GetProductVersion() string {
	return c.Version
}

func (c *KaiProductConfiguration) GetProductDescription() string {
	return c.Description
}

func (c *KaiProductConfiguration) UpdateProductVersion(version string, description ...string) error {
	if semver.IsValid(version) {
		c.Version = version

		if len(description) > 0 && description[0] != "" {
			c.Description = description[0]
		}

		return nil
	}

	return ErrInvalidVersion
}
