package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"gopkg.in/yaml.v3"
)

// Save writes config to file.
func (c *Config) Save() error {
	err := c.ensureConfigDir()
	if err != nil {
		return err
	}

	const writePermission = 0o600

	cfgFile, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, writePermission)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	_, err = cfgFile.Write(data)

	return err
}

func (c *Config) readFile() error {
	err := c.ensureConfigFile()
	if err != nil {
		return err
	}

	f, err := os.Open(c.Filename)
	if err != nil {
		return pathError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)

	return err
}

func (c *Config) ensureConfigDir() error {
	_, err := os.Stat(path.Dir(c.Filename))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	const writePermission = 0o771

	if os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(c.Filename), writePermission)
		if err != nil {
			return pathError(err)
		}
	}

	return nil
}

func (c *Config) ensureConfigFile() error {
	err := c.ensureConfigDir()
	if err != nil {
		return err
	}

	_, err = os.Stat(c.Filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		file, err := os.Create(c.Filename)
		if err != nil {
			return err
		}
		defer file.Close()

		data, err := yaml.Marshal(&c)
		if err != nil {
			return err
		}

		_, err = file.Write(data)

		return err
	}

	return nil
}

func pathError(err error) error {
	var pathError *os.PathError
	if errors.As(err, &pathError) && errors.Is(pathError.Err, syscall.ENOTDIR) {
		if p := findRegularFile(pathError.Path); p != "" {
			return fmt.Errorf("remove or rename regular file `%s` (must be a directory)", p) //nolint: goerr113
		}
	}

	return err
}

func findRegularFile(p string) string {
	for {
		if s, err := os.Stat(p); err == nil && s.Mode().IsRegular() {
			return p
		}

		newPath := path.Dir(p)
		if newPath == p || newPath == "/" || newPath == "." {
			break
		}

		p = newPath
	}

	return ""
}
