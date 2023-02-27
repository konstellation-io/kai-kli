package validate_test

import (
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kli/internal/krt/fileutils/compression"
	"github.com/konstellation-io/kli/internal/krt/validate"
	"github.com/konstellation-io/kli/mocks"
)

type testValidatorSuite struct {
	ctrl  *gomock.Controller
	mocks suiteMocks
}

type suiteMocks struct {
	logger   *mocks.MockLogger
	renderer *mocks.MockRenderer
}

func newTestValidatorSuite(t *testing.T) *testValidatorSuite {
	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	renderer := mocks.NewMockRenderer(ctrl)

	return &testValidatorSuite{
		ctrl,
		suiteMocks{
			logger:   logger,
			renderer: renderer,
		},
	}
}

func TestValidator_ValidateKrt(t *testing.T) {
	testCases := []struct {
		name                string
		testAsset           string
		wantValidationError bool
		wantExecutionError  bool
	}{
		{
			name:                "successfully validate a krt file",
			testAsset:           path.Join("testdata", "valid-krt.krt"),
			wantValidationError: false,
			wantExecutionError:  false,
		},
		{
			name:                "a krt is not valid if some node src doesn't exist",
			testAsset:           path.Join("testdata", "krt-invalid-src.krt"),
			wantValidationError: true,
			wantExecutionError:  true,
		},
		{
			name:                "a krt is not valid if entrypoint proto doesn't exist",
			testAsset:           path.Join("testdata", "krt-invalid-proto.krt"),
			wantValidationError: true,
			wantExecutionError:  true,
		},
		{
			name:                "a krt is not valid if krt.yaml is not valid",
			testAsset:           path.Join("testdata", "krt_invalid_yaml.krt"),
			wantValidationError: true,
			wantExecutionError:  true,
		},
		{
			name:                "a krt is not valid if it has no .krt extension",
			testAsset:           path.Join("testdata", "not-a-krt"),
			wantValidationError: false,
			wantExecutionError:  true,
		},
		{
			name:                "validation fails if krt file has no krt yaml",
			testAsset:           path.Join("testdata", "krt-invalid-no-yaml.krt"),
			wantValidationError: false,
			wantExecutionError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestValidatorSuite(t)
			compressor := compression.NewKrtCompressor()
			validator := validate.NewValidator(s.mocks.logger, compressor, s.mocks.renderer)

			if tc.wantValidationError {
				s.mocks.renderer.EXPECT().RenderValidationErrors(gomock.Len(1)).Return()
			} else {
				s.mocks.renderer.EXPECT().RenderValidationErrors(gomock.Nil()).Return().AnyTimes()
			}

			err := validator.ValidateKrt(tc.testAsset)
			if tc.wantExecutionError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
