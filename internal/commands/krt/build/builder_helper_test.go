package build_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kli/internal/testhelpers"
)

func generateDefaultKrtDir(t *testing.T, filePath, krtFileName string) {
	sampleKrtString := testhelpers.NewKrtBuilder().AsString()
	defaultTestKrtFiles := []string{
		"src/py-test/main.py",
		"src/py-test/main.pyc",
		"src/go-test/go.go",
		"docs/README.md",
	}

	createTestKrtContent(t, filePath, defaultTestKrtFiles...)
	createFile(t, filePath, sampleKrtString, krtFileName)
}

func createFile(t *testing.T, filePath, content, filename string) {
	t.Helper()

	err := os.WriteFile(filepath.Join(filePath, filename), []byte(content), 0600)
	assert.NoError(t, err)
}

func createTestKrtContent(t *testing.T, root string, files ...string) {
	t.Helper()

	for _, name := range files {
		name = path.Join(root, name)
		filePath := path.Dir(name)
		_, err := os.Stat(filePath)

		if os.IsNotExist(err) {
			err = os.MkdirAll(filePath, 0755)
			assert.NoError(t, err)
		}

		f, err := os.Create(name)

		defer func() { //nolint:gocritic
			_ = f.Close()
		}()
		assert.NoError(t, err)
	}
}
