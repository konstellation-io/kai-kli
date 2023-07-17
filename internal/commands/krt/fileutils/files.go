package fileutils

import (
	"errors"
	"os"
	"path/filepath"
)

// GetKrtYaml returns the full path, and the file name.
func GetKrtYaml(src string) (yamlPath, yamlFormat string, err error) {
	yamlFormats := []string{"krt.yaml", "krt.yml"}

	for _, yamlFormat := range yamlFormats {
		yamlPath := filepath.Join(src, yamlFormat)

		if fileExists(yamlPath) {
			return yamlPath, yamlFormat, nil
		}
	}

	return "", "", errors.New("no yaml file found")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Copy(source, target string) error {
	bytesRead, err := os.ReadFile(source)
	if err != nil {
		return errors.New("error reading file")
	}

	err = os.WriteFile(target, bytesRead, 0664) //nolint:gosec,gomnd
	if err != nil {
		return errors.New("error writing file")
	}

	return nil
}
