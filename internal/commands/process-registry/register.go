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
	ErrZipFileCouldNotBeCreated = errors.New("zip file could not be created")
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

	tmpZipFile, err := c.creteTempZipFile(sourcesPath, dockerfilePath)
	if err != nil {
		return ErrZipFileCouldNotBeCreated
	}

	registeredProcess, err := c.processRegistryClient.
		Register(srv, tmpZipFile, productID, processID, processType, version)
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

	writer := gzip.NewWriter(f)
	defer writer.Close()

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	// 2. Go through all the files of the source
	for _, path := range paths {
		if err := c.addToZipFile(writer, tarWriter, path); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (c *Handler) addToZipFile(writer *gzip.Writer, tarWriter *tar.Writer, sourcePath string) error {
	return filepath.WalkDir(sourcePath, func(filePath string, info os.DirEntry, err error) error {
		c.logger.Debug(fmt.Sprintf("Adding %s to zip file, error %s\n", filePath, err))
		if err != nil {
			return err
		}
		if filePath == sourcePath {
			return nil
		}

		fileInfo, err := info.Info()
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(sourcePath, filePath)
		if err != nil {
			return err
		}

		// generate tar header
		header, err := tar.FileInfoHeader(fileInfo, relativePath)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relativePath)

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		// if not a dir, write file content
		if !fileInfo.IsDir() {
			data, err := os.Open(filePath)
			if err != nil {
				return err
			}

			_, err = io.Copy(tarWriter, data)
			if err != nil {
				return err
			}
		}

		return nil
		//
		//// 3. Create a local file header
		//header, err := tar.FileInfoHeader(info, relativePath)
		//if err != nil {
		//	return err
		//}
		//
		//// set compression
		////header.Method = zip.Deflate
		//
		//// 4. Set relative filePath of a file as the header name
		//header.Name, err = filepath.Rel(filepath.Dir(sourcePath), filePath)
		//if err != nil {
		//	return err
		//}
		//header.Name = filepath.ToSlash(relativePath)
		//
		//// 5. Create writer for the file header and save content of the file
		//err = tarWriter.WriteHeader(header)
		//if err != nil {
		//	return err
		//}
		//
		//if !file.IsDir() {
		//	data, err := os.Open(filePath)
		//	if err != nil {
		//		return err
		//	}
		//
		//	_, err = io.Copy(c.tarWriter, data)
		//	if err != nil {
		//		return err
		//	}
		//}

		//	return nil
		//}
		//
		//	f, err := os.Open(filePath)
		//	if err != nil {
		//		return err
		//	}
		//	defer f.Close()
		//
		//	_, err = io.Copy(headerWriter, f)
		//	return err
	})
}

func (c *Handler) pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
