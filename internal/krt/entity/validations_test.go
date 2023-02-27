package entity_test

import (
	"testing"

	"github.com/konstellation-io/kli/internal/krt/entity"
	"github.com/konstellation-io/kli/internal/krt/errors"
	"github.com/stretchr/testify/assert"
)

func TestFile_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		krt           *entity.File
		isValid       bool
		invalidFields []string
	}{
		{
			name:          "valid krt validation returns no errors",
			krt:           NewKrtBuilder().Build(),
			isValid:       true,
			invalidFields: nil,
		},
		{
			name:          "valid krt validation returns errors",
			krt:           NewKrtBuilder().WithKrtVersion("v1").Build(),
			isValid:       false,
			invalidFields: []string{"krtVersion"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errs := tc.krt.Validate()
			if !tc.isValid {
				assert.ElementsMatchf(t,
					tc.invalidFields,
					getInvalidFieldsFromErrors(errs),
					"expected: %v, actual: %v", tc.invalidFields, getInvalidFieldsFromErrors(errs),
				)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func getInvalidFieldsFromErrors(errs []*errors.ValidationError) []string {
	invalidFields := make([]string, 0, len(errs))
	for _, err := range errs {
		invalidFields = append(invalidFields, err.Field)
	}

	return invalidFields
}
