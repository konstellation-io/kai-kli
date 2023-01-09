package builder

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/internal/krt"
	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/mattn/go-zglob"
)

var ErrInvalidVersionName = errors.New("invalid version name")

type Builder struct {
	logger    simplelogger.SimpleLoggerInterface
	validator *krt.CustomValidator
}

// New creates a new Builder instance.
func New(validator *krt.CustomValidator) Builder {
	return NewWithLogger(simplelogger.New(simplelogger.LevelInfo), validator)
}

// NewWithLogger creates a new Builder instance with a custom logger.
func NewWithLogger(l simplelogger.SimpleLoggerInterface, validator *krt.CustomValidator) Builder {
	return Builder{logger: l, validator: validator}
}

// cleanPatterns puts a set of given patterns into their context.
func (b *Builder) cleanPatterns(src string, patterns []string) []string {
	for _, p := range patterns {
		p = filepath.ToSlash(p)
		t := filepath.Join(src, p)
		info, err := os.Stat(t)

		if err != nil {
			if strings.HasSuffix(p, "*") {
				patterns = append(patterns, fmt.Sprintf("%s*/*", p))
			} else {
				patterns = append(patterns, p)
			}

			continue
		}

		if info.IsDir() {
			patterns = append(patterns, fmt.Sprintf("%s**/*", p))
		}
	}

	return patterns
}

// getIgnorePatterns returns clean ignore patterns for a given directory.
func (b *Builder) getIgnorePatterns(dir string) []string {
	ignoreFile := filepath.Join(dir, ".krtignore")

	data, err := ioutil.ReadFile(ignoreFile) //nolint:gocritic
	if err != nil {
		b.logger.Infof("Skipping .krtignore file.")
		return nil
	}

	p := strings.Split(string(data), "\n")
	patterns := b.cleanPatterns(dir, p)

	return patterns
}

// skipFile checks with a set of clean patterns if a file has to be skipped or not.
func (b *Builder) skipFile(file string, patterns []string) (bool, error) {
	if file == "" {
		return false, nil
	}

	for _, p := range patterns {
		if p == file {
			return true, nil
		}

		match, err := zglob.Match(p, file)

		if match || err != nil {
			return match, err
		}
	}

	return false, nil
}

// CreateKrt create a krt file from a source dir.
func (b *Builder) CreateKrt(src, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}

	defer zipfile.Close()

	zr, err := gzip.NewWriterLevel(zipfile, gzip.BestCompression)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(zr)

	patterns := b.getIgnorePatterns(src)

	err = filepath.Walk(src, func(path string, info os.FileInfo, _ error) error {
		if path == src {
			return nil
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		matchName := strings.Replace(strings.TrimPrefix(path, src), "/", "", 1)
		skip, err := b.skipFile(matchName, patterns)
		if err != nil {
			return err
		}

		if skip {
			return nil
		}

		// generate tar header
		header, err := tar.FileInfoHeader(info, relativePath)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relativePath)
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		// if not a dir, write file content
		if !info.IsDir() {
			data, err := os.Open(path)
			if err != nil {
				return err
			}

			_, err = io.Copy(tw, data)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}

	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}

	return nil
}

// Build creates a KrtFile.
func (b *Builder) Build(src, target string) error {
	err := b.CreateKrt(src, target)
	if err != nil {
		return err
	}

	return nil
}

// UpdateVersion updates krt yaml file version name.
func (b *Builder) UpdateVersion(krtFile *krt.File, yamlPath, version string) error {
	const writePermission = 0664

	if version == "" {
		return ErrInvalidVersionName
	}

	validVersion := b.validator.ValidateVersionName(version)

	if !validVersion {
		return ErrInvalidVersionName
	}

	krtFile.Version = version

	file, err := os.OpenFile(yamlPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, writePermission)

	if err != nil {
		return fmt.Errorf("error while opening yaml file: %w", err)
	}
	defer file.Close()

	data, err := yaml.Marshal(&krtFile)
	if err != nil {
		return fmt.Errorf("error while marshaling yaml file: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error while writing yaml file: %w", err)
	}

	return nil
}
