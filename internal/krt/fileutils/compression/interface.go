package compression

import (
	"os"
)

//go:generate mockgen -source=${GOFILE} -destination=../../../../mocks/compressor.go -package=mocks

type Compressor interface {
	NewCompressedFile(filePath string) (CompressedFile, error)
	Extract(krtPath, dst string) error
}

type CompressedFile interface {
	AddFile(file os.FileInfo, path, relativePath string) error
	Compress() error
}
