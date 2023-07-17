package build_test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	matchers "github.com/Storytel/gomock-matchers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kli/internal/commands/krt/build"
	"github.com/konstellation-io/kli/internal/commands/krt/entity"
	"github.com/konstellation-io/kli/mocks"
)

const builderTmpPath = "krt-build-*"

type testBuilderSuite struct {
	ctrl  *gomock.Controller
	mocks suiteMocks
}

type suiteMocks struct {
	logger     *mocks.MockLogger
	compressor *mocks.MockCompressor
}

func newTestBuilderSuite(t *testing.T) *testBuilderSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	compressor := mocks.NewMockCompressor(ctrl)
	mocks.AddLoggerExpects(logger)

	return &testBuilderSuite{
		ctrl: ctrl,
		mocks: suiteMocks{
			logger:     logger,
			compressor: compressor,
		},
	}
}

func TestBuilder_Build(t *testing.T) {
	s := newTestBuilderSuite(t)

	// GIVEN a local directory
	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s-*", t.Name()))
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	// and a krt file
	krtFileName := "krt.yaml"
	// and that the local directory has the content of a valid krt
	generateDefaultKrtDir(t, tmp, krtFileName)

	// WHEN krt is build with an empty version and disabling the updateLocal option
	buildDstPath := "."
	buildVersion := ""
	updateLocal := false

	// THEN a compressed krt file named with the version in the krt.yaml is created
	compressedFile := mocks.NewMockCompressedFile(s.ctrl)
	s.mocks.compressor.EXPECT().NewCompressedFile("version-name.krt").Return(compressedFile, nil)

	// and all files in the krt directory are included in alphabetical order starting with the files
	gomock.InOrder(
		// and the krt included in the compressed file a temporal one created by the builder
		compressedFile.EXPECT().AddFile(gomock.Any(), matchers.Regexp(builderTmpPath), krtFileName).Return(nil),
		compressedFile.EXPECT().AddFile(gomock.Any(), path.Join(tmp, "docs"), "docs").Return(nil),
		compressedFile.EXPECT().AddFile(gomock.Any(), path.Join(tmp, "docs/README.md"), "docs/README.md").Return(nil),
		compressedFile.EXPECT().AddFile(gomock.Any(), path.Join(tmp, "src"), "src").Return(nil),
		compressedFile.EXPECT().AddFile(gomock.Any(), path.Join(tmp, "src/go-test"), "src/go-test").Return(nil),
		compressedFile.EXPECT().AddFile(gomock.Any(), path.Join(tmp, "src/go-test/go.go"), "src/go-test/go.go").Return(nil),
		compressedFile.EXPECT().AddFile(gomock.Any(), path.Join(tmp, "src/py-test"), "src/py-test").Return(nil),
		compressedFile.EXPECT().
			AddFile(gomock.Any(), path.Join(tmp, "src/py-test/main.py"), "src/py-test/main.py").Return(nil),
		compressedFile.EXPECT().
			AddFile(gomock.Any(), path.Join(tmp, "src/py-test/main.pyc"), "src/py-test/main.pyc").Return(nil),
	)

	compressedFile.EXPECT().Compress().Return(nil)

	b := build.NewBuilder(s.mocks.logger, s.mocks.compressor)
	err = b.Build(tmp, buildDstPath, buildVersion, updateLocal)
	assert.NoError(t, err)
}

func TestBuilder_BuildWithIgnoredFiles(t *testing.T) {
	s := newTestBuilderSuite(t)

	// Create test dir structure
	tmp, err := os.MkdirTemp("", "TestBuilder_*")
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	krtFileName := "krt.yaml"
	generateDefaultKrtDir(t, tmp, krtFileName)
	createFile(t, tmp, "*.pyc", ".krtignore")

	compressedFile := mocks.NewMockCompressedFile(s.ctrl)
	s.mocks.compressor.EXPECT().NewCompressedFile("version-name.krt").Return(compressedFile, nil)

	compressedFile.EXPECT().
		AddFile(gomock.Any(), gomock.Not("src/py-test/main.pyc"), gomock.Not("src/py-test/main.pyc")).Return(nil).AnyTimes()

	compressedFile.EXPECT().Compress().Return(nil)

	b := build.NewBuilder(s.mocks.logger, s.mocks.compressor)
	err = b.Build(tmp, ".", "", false)
	assert.NoError(t, err)
}

func TestBuilder_Flags(t *testing.T) {
	testCases := []struct {
		name         string
		versionName  string
		yamlFileName string
		updateLocal  bool
		wantError    bool
	}{
		{
			name:         "set version without updating local yaml",
			versionName:  "new-version",
			yamlFileName: "krt.yaml",
			updateLocal:  false,
			wantError:    false,
		},
		{
			name:         "set version updating local yaml",
			versionName:  "new-version",
			yamlFileName: "krt.yaml",
			updateLocal:  true,
			wantError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestBuilderSuite(t)
			// GIVEN that a temporal path is created
			tmp, err := os.MkdirTemp("", "TestBuilder_*")
			assert.NoError(t, err)
			defer os.RemoveAll(tmp)

			// and the temporal path contains a default krt structure
			generateDefaultKrtDir(t, tmp, tc.yamlFileName)

			// WHEN and user builds a krt with the test case flags
			// THEN a new compressed file is generated
			compressedKrt := fmt.Sprintf("%s.krt", tc.versionName)
			compressedFile := mocks.NewMockCompressedFile(s.ctrl)
			s.mocks.compressor.EXPECT().NewCompressedFile(compressedKrt).Return(compressedFile, nil)

			// and the local krt.yaml file is included in the compressed krt
			localKrtYamlPath := path.Join(tmp, tc.yamlFileName)
			if tc.updateLocal {
				compressedFile.
					EXPECT().
					AddFile(gomock.Any(), localKrtYamlPath, tc.yamlFileName).
					Return(nil).
					Times(1)
			} else {
				compressedFile.
					EXPECT().
					AddFile(gomock.Any(), matchers.Regexp(builderTmpPath), tc.yamlFileName).
					Return(nil).
					Times(1)
			}

			// and the rest of the files in the directory are included
			compressedFile.
				EXPECT().
				AddFile(gomock.Any(), gomock.Not(localKrtYamlPath), gomock.Not(tc.yamlFileName)).
				Return(nil).
				AnyTimes()

			// and the compressed file is closed generating the output
			compressedFile.EXPECT().Compress().Return(nil)

			b := build.NewBuilder(s.mocks.logger, s.mocks.compressor)
			err = b.Build(tmp, ".", tc.versionName, tc.updateLocal)
			// and there is no errors
			if tc.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tc.updateLocal {
				// and the version name is updated in the local krt.yaml file
				modifiedKrt, err := entity.ParseFile(localKrtYamlPath)
				assert.NoError(t, err)
				assert.Equal(t, tc.versionName, string(modifiedKrt.Version))
			}
		})
	}
}

func TestBuilder_ErrorIfLocalYamlNotExist(t *testing.T) {
	s := newTestBuilderSuite(t)
	tmp, err := os.MkdirTemp("", "TestBuilder_*")
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	newVersionName := "new-version"
	krtFileName := "krt.yaml"

	generateDefaultKrtDir(t, tmp, krtFileName)
	err = os.Remove(path.Join(tmp, krtFileName))
	assert.NoError(t, err)

	b := build.NewBuilder(s.mocks.logger, s.mocks.compressor)
	err = b.Build(tmp, ".", newVersionName, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error getting krt yaml path")
}

func TestBuilder_ErrorIfLocalYamlHasInvalidContent(t *testing.T) {
	s := newTestBuilderSuite(t)
	tmp, err := os.MkdirTemp("", "TestBuilder_*")
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	newVersionName := "new-version"
	krtFileName := "krt.yaml"

	generateDefaultKrtDir(t, tmp, krtFileName)
	err = os.Remove(path.Join(tmp, krtFileName))
	assert.NoError(t, err)
	createFile(t, tmp, "this is not a krt structure", krtFileName)

	b := build.NewBuilder(s.mocks.logger, s.mocks.compressor)
	err = b.Build(tmp, ".", newVersionName, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error parsing krt yaml")
}

func TestBuilder_ErrorIfLocalKrtYamlCouldNotBeCompressed(t *testing.T) {
	s := newTestBuilderSuite(t)
	tmp, err := os.MkdirTemp("", "TestBuilder_*")
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	versionName := "version-name"
	krtFileName := "krt.yaml"

	generateDefaultKrtDir(t, tmp, krtFileName)
	assert.NoError(t, err)

	compressedKrt := fmt.Sprintf("%s.krt", versionName)
	compressedFile := mocks.NewMockCompressedFile(s.ctrl)
	s.mocks.compressor.EXPECT().NewCompressedFile(compressedKrt).Return(compressedFile, nil)
	compressedFile.EXPECT().AddFile(gomock.Any(), gomock.Any(), krtFileName).Return(errors.New("error compressing file"))

	b := build.NewBuilder(s.mocks.logger, s.mocks.compressor)
	err = b.Build(tmp, ".", versionName, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error compressing file")
}

func TestBuilder_ErrorIfSomeFileCouldNotBeCompressed(t *testing.T) {
	s := newTestBuilderSuite(t)
	tmp, err := os.MkdirTemp("", "TestBuilder_*")
	assert.NoError(t, err)

	defer os.RemoveAll(tmp)

	versionName := "version-name"
	krtFileName := "krt.yaml"

	generateDefaultKrtDir(t, tmp, krtFileName)
	assert.NoError(t, err)

	compressedKrt := fmt.Sprintf("%s.krt", versionName)
	compressedFile := mocks.NewMockCompressedFile(s.ctrl)
	s.mocks.compressor.EXPECT().NewCompressedFile(compressedKrt).Return(compressedFile, nil)
	compressedFile.EXPECT().AddFile(gomock.Any(), gomock.Any(), krtFileName).Return(nil)
	compressedFile.EXPECT().AddFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error compressing file"))

	b := build.NewBuilder(s.mocks.logger, s.mocks.compressor)
	err = b.Build(tmp, ".", versionName, false)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error compressing file")
}
