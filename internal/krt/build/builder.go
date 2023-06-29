package build

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"

	"github.com/konstellation-io/kli/internal/krt/entity"
	"github.com/konstellation-io/kli/internal/krt/fileutils"
	"github.com/konstellation-io/kli/internal/krt/fileutils/compression"
	"github.com/konstellation-io/kli/internal/logging"
)

type Builder struct {
	logger     logging.Interface
	compressor compression.Compressor
}

// NewBuilder creates a new Builder instance.
func NewBuilder(logger logging.Interface, compressor compression.Compressor) *Builder {
	return &Builder{
		logger:     logger,
		compressor: compressor,
	}
}

func (b *Builder) Build(krtSrcPath, krtDstPath, version string, updateLocal bool) error {
	yamlPath, yamlBaseName, err := fileutils.GetKrtYaml(krtSrcPath)
	if err != nil {
		return fmt.Errorf("error getting krt yaml path: %w", err)
	}

	localConfig, err := entity.ParseFile(yamlPath)
	if err != nil {
		return fmt.Errorf("error parsing krt yaml: %w", err)
	}

	buildVersion := string(localConfig.Version)

	if version != "" {
		buildVersion = version
	}

	if !updateLocal {
		tmpDir, err := os.MkdirTemp("", "krt-build-*")
		if err != nil {
			return fmt.Errorf("error creating temporal directory: %w", err)
		}

		defer os.RemoveAll(tmpDir)

		tmpYamlPath := path.Join(tmpDir, yamlBaseName)

		err = fileutils.Copy(yamlPath, tmpYamlPath)
		if err != nil {
			return fmt.Errorf("error copying local krt yaml to temporal directory: %w", err)
		}

		yamlPath = tmpYamlPath
	}

	err = fileutils.UpdateVersion(yamlPath, yamlPath, buildVersion)
	if err != nil {
		return fmt.Errorf("error updating version: %w", err)
	}

	krtFilePath := path.Join(krtDstPath, fmt.Sprintf("%s.krt", buildVersion))

	err = b.buildKrt(krtSrcPath, krtFilePath, yamlPath)
	if err != nil {
		return fmt.Errorf("error building krt file: %w", err)
	}

	b.logger.Success("New KRT file created.")

	return nil
}

// buildKrt builds a krt file from a source dir.
func (b *Builder) buildKrt(src, target, yamlPath string) error {
	compressedFile, err := b.compressor.NewCompressedFile(target)
	if err != nil {
		return err
	}

	krtignorePath := path.Join(src, ".krtignore")
	ignoreLocalYaml := yamlPath != ""
	patterns := b.getIgnorePatterns(krtignorePath, ignoreLocalYaml)

	if ignoreLocalYaml {
		fileInfo, err := os.Stat(yamlPath)
		if err != nil {
			return err
		}

		relativePath := path.Base(yamlPath)

		err = compressedFile.AddFile(fileInfo, yamlPath, relativePath)
		if err != nil {
			return err
		}
	}

	err = b.compressFiles(src, patterns, compressedFile)
	if err != nil {
		return err
	}

	return compressedFile.Compress()
}

func (b *Builder) compressFiles(src string, patterns ignore.IgnoreParser, compressedFile compression.CompressedFile) error {
	return filepath.WalkDir(src, func(path string, info os.DirEntry, _ error) error {
		// Ignore root dir folder
		if path == src {
			return nil
		}

		matchName := strings.Replace(strings.TrimPrefix(path, src), "/", "", 1)
		skip := patterns.MatchesPath(matchName)

		if skip {
			b.logger.Debug(fmt.Sprintf("Skipped file: %s", path))
			return nil
		}

		fileInfo, err := info.Info()
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		err = compressedFile.AddFile(fileInfo, path, relativePath)
		if err != nil {
			return err
		}

		return nil
	})
}

func (b *Builder) getIgnorePatterns(krtignorePath string, ignoreLocalYaml bool) ignore.IgnoreParser {
	var additionalLines []string
	if ignoreLocalYaml {
		additionalLines = append(additionalLines, "krt.yml", "krt.yaml")
	}

	patterns, err := ignore.CompileIgnoreFileAndLines(krtignorePath, additionalLines...)
	if err != nil {
		b.logger.Debug(fmt.Sprintf("Error reading .krtignore file: %s", err))
		b.logger.Info("Ignoring .krtignore file")

		return ignore.CompileIgnoreLines(additionalLines...)
	}

	return patterns
}
