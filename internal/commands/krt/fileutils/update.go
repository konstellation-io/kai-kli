package fileutils

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/konstellation-io/kli/internal/commands/krt/entity"
)

func UpdateVersion(yamlPath, outputFile, version string) error {
	_, err := entity.NewResourceName(version)
	if err != nil {
		return errors.New("invalid version name")
	}

	const writePermission = 0664

	actualKrt, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("error reading krt yaml file: %w", err)
	}

	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, writePermission)
	if err != nil {
		return fmt.Errorf("error while opening yaml file: %w", err)
	}
	defer file.Close()

	actualKrtStr := string(actualKrt)
	versionLine := regexp.MustCompile("version:.*\n?")

	var updatedKrt string
	if !versionLine.MatchString(actualKrtStr) {
		updatedKrt = fmt.Sprintf("version: %s\n%s", version, actualKrtStr)
	} else {
		updatedKrt = versionLine.ReplaceAllString(actualKrtStr, fmt.Sprintf("version: %s\n", version))
	}

	_, err = file.WriteString(updatedKrt)
	if err != nil {
		return fmt.Errorf("error writing yaml file: %w", err)
	}

	return nil
}
