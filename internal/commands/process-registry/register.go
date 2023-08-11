package processregistry

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/konstellation-io/krt/pkg/krt"
)

var (
	ErrPathDoesNotExist         = errors.New("path does not exist")
	ErrZipFileCouldNotBeCreated = errors.New("zip file could not be created")
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

	tmpZipFile, err := c.creteTempZipFile(opts.SourcesPath, opts.Dockerfile)
	if err != nil {
		return ErrZipFileCouldNotBeCreated
	}

	registeredProcess, err := c.processRegistryClient.
		Register(srv, tmpZipFile, opts.ProductID, opts.ProcessID, string(opts.ProcessType), opts.Version)
	if err != nil {
		return err
	}

	c.logger.Success(fmt.Sprintf("Process registered with ID: %s", registeredProcess))

	return nil
}

func (c *Handler) creteTempZipFile(paths ...string) (*os.File, error) {
	tmpPath := os.TempDir()

	f, err := os.CreateTemp(tmpPath, "process-*.zip")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	for _, path := range paths {
		if err := c.addToZipFile(writer, path); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (c *Handler) addToZipFile(writer *zip.Writer, sourcePath string) error {
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		c.logger.Debug(fmt.Sprintf("Adding %s to zip file, error %s\n", path, err))
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(sourcePath), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
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

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func (c *Handler) pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
