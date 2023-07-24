package processregistry

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrPathDoesNotExist         = errors.New("path does not exist")
	ErrZipFileCouldNotBeCreated = errors.New("tar.gz file could not be created")
)

func (c *Handler) RegisterProcess(serverName, productID, processType, processID,
	sourcesPath, dockerfilePath, version string) error {
	kaiConfig, err := c.configService.GetConfiguration()
	if err != nil {
		return err
	}

	srv, err := kaiConfig.GetServerOrDefault(serverName)
	if err != nil {
		return err
	}

	if !c.pathExists(sourcesPath) || !c.pathExists(dockerfilePath) {
		return ErrPathDoesNotExist
	}

	tmpZipFile, err := c.createTempTarGzFile(sourcesPath, dockerfilePath)
	if err != nil {
		return ErrZipFileCouldNotBeCreated
	}

	defer tmpZipFile.Close()

	registeredProcess, err := c.processRegistryClient.
		Register(srv, tmpZipFile, productID, processID, processType, version)
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
