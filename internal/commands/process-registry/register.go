package processregistry

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/konstellation-io/krt/pkg/krt"
)

var (
	ErrPathDoesNotExist         = errors.New("path does not exist")
	ErrZipFileCouldNotBeCreated = errors.New("tar.gz file could not be created")
)

type RegisterProcessOpts struct {
	ServerName  string
	ProductID   string
	ProcessType krt.ProcessType
	ProcessID   string
	SourcesPath string
	Dockerfile  string
	Version     string
}

func (c *Handler) RegisterProcess(opts *RegisterProcessOpts) error {
	if !opts.ProcessType.IsValid() {
		return fmt.Errorf("invalid process type: %q", opts.ProcessType)
	}

	kaiConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(opts.ServerName)
	if err != nil {
		return err
	}

	if !c.pathExists(opts.SourcesPath) || !c.pathExists(opts.Dockerfile) {
		return ErrPathDoesNotExist
	}

	if info, err := os.Stat(opts.SourcesPath); err == nil && info.IsDir() {
		if !strings.HasSuffix(opts.SourcesPath, "/") {
			opts.SourcesPath += "/"
		}
	}

	tmpZipFile, err := c.createTempTarGzFile(opts.SourcesPath, opts.Dockerfile)
	if err != nil {
		return ErrZipFileCouldNotBeCreated
	}

	defer tmpZipFile.Close()

	registeredProcess, err := c.processRegistryClient.
		Register(srv, tmpZipFile, opts.ProductID, opts.ProcessID, string(opts.ProcessType), opts.Version)
	if err != nil {
		return err
	}

	c.logger.Success(fmt.Sprintf("Process registered with ID: %s", registeredProcess))

	return nil
}

func (c *Handler) createTempTarGzFile(paths ...string) (*os.File, error) {
	tmpPath := os.TempDir()

	f, err := os.CreateTemp(tmpPath, "process-*.tar.gz")
	if err != nil {
		return nil, err
	}

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, path := range paths {
		if err := c.addToTarGz(tw, path); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (c *Handler) addToTarGz(tw *tar.Writer, sourcePath string) error {
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		c.logger.Debug(fmt.Sprintf("Adding %s to tar.gz file, error %s\n", path, err))
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(filepath.Dir(sourcePath), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})
}

func (c *Handler) pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
