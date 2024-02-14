package processregistry

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
	"github.com/konstellation-io/krt/pkg/krt"
	ignore "github.com/sabhiram/go-gitignore"
)

const (
	_krtignoreFileName = ".krtignore"
)

var (
	ErrPathDoesNotExist                 = errors.New("path does not exist")
	ErrZipFileCouldNotBeCreated         = errors.New("tar.gz file could not be created")
	ErrInvalidProcessType               = errors.New("invalid process type")
	ErrMissingProduct                   = errors.New("missing product")
	ErrIncompatibleProductAndPublicOpts = errors.New("product and public args are incompatible")
)

type RegisterProcessOpts struct {
	ServerName  string
	ProductID   string
	ProcessType krt.ProcessType
	ProcessID   string
	SourcesPath string
	Dockerfile  string
	Version     string
	IsPublic    bool
}

func (o *RegisterProcessOpts) Validate() error {
	if !o.ProcessType.IsValid() {
		return ErrInvalidProcessType
	}

	if o.ProductID == "" && !o.IsPublic {
		return ErrMissingProduct
	}

	if o.ProductID != "" && o.IsPublic {
		return ErrIncompatibleProductAndPublicOpts
	}

	return nil
}

func (c *Handler) RegisterProcess(opts *RegisterProcessOpts) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validating args: %w", err)
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

	patternsToIgnore := c.getPattersToIgnore(opts.SourcesPath)

	if info, err := os.Stat(opts.SourcesPath); err == nil && info.IsDir() {
		if !strings.HasSuffix(opts.SourcesPath, "/") {
			opts.SourcesPath += "/"
		}
	}

	tmpZipFile, err := c.createTempTarGzFile(patternsToIgnore, opts.SourcesPath, opts.Dockerfile)
	if err != nil {
		return ErrZipFileCouldNotBeCreated
	}

	defer tmpZipFile.Close()

	registeredProcess, err := c.registerProcess(srv, tmpZipFile, opts)
	if err != nil {
		return err
	}

	c.renderer.RenderProcessRegistered(registeredProcess)

	return nil
}

func (c *Handler) registerProcess(srv *configuration.Server, tmpZipFile *os.File,
	opts *RegisterProcessOpts) (*entity.RegisteredProcess, error) {
	if opts.IsPublic {
		return c.processRegistryClient.RegisterPublic(srv, tmpZipFile, opts.ProcessID, string(opts.ProcessType), opts.Version)
	}

	return c.processRegistryClient.
		Register(srv, tmpZipFile, opts.ProductID, opts.ProcessID, string(opts.ProcessType), opts.Version)
}

func (c *Handler) createTempTarGzFile(patternsToIgnore ignore.IgnoreParser, paths ...string) (*os.File, error) {
	tmpPath := os.TempDir()

	f, err := os.CreateTemp(tmpPath, "process-*.tar.gz")
	if err != nil {
		return nil, err
	}

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, dirPath := range paths {
		if err := c.addToTarGz(tw, dirPath, patternsToIgnore); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (c *Handler) addToTarGz(tw *tar.Writer, sourcePath string, patternsToIgnore ignore.IgnoreParser) error {
	return filepath.Walk(sourcePath, func(dirPath string, info os.FileInfo, err error) error {
		if patternsToIgnore.MatchesPath(path.Base(dirPath)) {
			c.logger.Debug(fmt.Sprintf("Skipping file %q", dirPath))
			return nil
		}

		c.logger.Debug(fmt.Sprintf("Adding %s to tar.gz file, error %s\n", dirPath, err))

		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(filepath.Dir(sourcePath), dirPath)
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

		f, err := os.Open(dirPath)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)

		return err
	})
}

func (c *Handler) pathExists(pathToCheck string) bool {
	if _, err := os.Stat(pathToCheck); os.IsNotExist(err) {
		return false
	}

	return true
}

func (c *Handler) getPattersToIgnore(dirPath string) ignore.IgnoreParser {
	krtignorePath := path.Join(dirPath, _krtignoreFileName)

	patterns, err := ignore.CompileIgnoreFile(krtignorePath)
	if err != nil {
		c.logger.Info("Ignoring .krtignore file")

		return ignore.CompileIgnoreLines()
	}

	return patterns
}
