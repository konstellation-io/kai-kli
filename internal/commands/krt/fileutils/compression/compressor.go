package compression

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type KrtCompressor struct{}

type KrtCompressedFile struct {
	file      *os.File
	tarWriter *tar.Writer
	zipWriter *gzip.Writer
	filePath  string
}

func NewKrtCompressor() *KrtCompressor {
	return &KrtCompressor{}
}

func (c *KrtCompressor) NewCompressedFile(filePath string) (CompressedFile, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	zipWriter, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	tarWriter := tar.NewWriter(zipWriter)

	return &KrtCompressedFile{
		file:      file,
		tarWriter: tarWriter,
		zipWriter: zipWriter,
		filePath:  filePath,
	}, nil
}

func (c *KrtCompressor) Extract(krtPath, dst string) error {
	// decompress file on temporal dir
	file, err := os.Open(krtPath)
	if err != nil {
		return fmt.Errorf("error opening krt file: %w", err)
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading header: %w", err)
		}

		dst := path.Join(dst, header.Name) //nolint: gosec

		switch header.Typeflag {
		case tar.TypeDir:
			permissions := 0755
			if err := os.Mkdir(dst, os.FileMode(permissions)); err != nil {
				return fmt.Errorf("error extracting directory: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(dst)
			if err != nil {
				return fmt.Errorf("error creating uncompressed file: %w", err)
			}

			defer outFile.Close() //nolint:gocritic

			if _, err := io.Copy(outFile, tarReader); err != nil { //nolint: gosec
				return fmt.Errorf("error extracting file: %w", err)
			}

		default:
			return fmt.Errorf("found a not supported compressed file")
		}
	}

	return nil
}

func (c *KrtCompressedFile) AddFile(file os.FileInfo, filePath, relativePath string) error {
	// generate tar header
	header, err := tar.FileInfoHeader(file, relativePath)
	if err != nil {
		return err
	}

	header.Name = filepath.ToSlash(relativePath)

	err = c.tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	// if not a dir, write file content
	if !file.IsDir() {
		data, err := os.Open(filePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(c.tarWriter, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *KrtCompressedFile) Compress() error {
	if err := c.tarWriter.Close(); err != nil {
		return err
	}

	if err := c.zipWriter.Close(); err != nil {
		return err
	}

	return c.file.Close()
}
