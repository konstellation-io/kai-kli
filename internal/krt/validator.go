package krt

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/konstellation-io/kli/internal/logger"
)

type CustomValidator struct {
	logger          logger.Interface
	fieldsValidator *validator.Validate
}

func NewCustomValidator(loggerI logger.Interface) *CustomValidator {
	return &CustomValidator{
		logger:          loggerI,
		fieldsValidator: initializeFieldsValidator(),
	}
}

func initializeFieldsValidator() *validator.Validate {
	krtValidator := validator.New()

	// register validator for resource names. Ex: name-valid123
	reResourceName := regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	_ = krtValidator.RegisterValidation("resource-name", func(fl validator.FieldLevel) bool {
		return reResourceName.MatchString(fl.Field().String())
	})

	// register validator for env var names. Ex: NAME_VALID123
	reEnvVar := regexp.MustCompile("^[A-Z0-9]([_A-Z0-9]*[A-Z0-9])?$")
	_ = krtValidator.RegisterValidation("env", func(fl validator.FieldLevel) bool {
		return reEnvVar.MatchString(fl.Field().String())
	})

	return krtValidator
}

func (v *CustomValidator) Run(krtFile *File) error {
	var errs []error

	v.logger.Info("Validating krt.yml")

	fieldValidationErrors := v.validateYamlFields(krtFile)

	if fieldValidationErrors != nil {
		errs = append(errs, fieldValidationErrors...)
	}

	workflowValidationErrors := v.getWorkflowsValidationErrors(krtFile.Workflows)
	if workflowValidationErrors != nil {
		errs = append(errs, workflowValidationErrors...)
	}

	if errs != nil {
		return NewValidationError(errs)
	}

	return nil
}

func (v *CustomValidator) validateYamlFields(yaml interface{}) []error {
	err := v.fieldsValidator.Struct(yaml)
	return v.getYamlFieldsValidationErrors(err)
}

//nolint:goerr113
func (v *CustomValidator) getYamlFieldsValidationErrors(err error) []error {
	if err == nil {
		return nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []error

		hasResNameErr := false

		for _, e := range errs {
			location := strings.Replace(e.Namespace(), "Krt.", "", 1)

			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Errorf("the field %q is required",
					location))
			case "lt":
				errorMessages = append(errorMessages, fmt.Errorf("invalid length %q at %q must be lower than %s",
					e.Value(), location, e.Param()))
			case "lte":
				errorMessages = append(errorMessages, fmt.Errorf("invalid length %q at %q must be lower or equal than %s",
					e.Value(), location, e.Param()))
			case "gt":
				errorMessages = append(errorMessages, fmt.Errorf("invalid length %q at %q must be greater than %s",
					e.Value(), location, e.Param()))
			case "gte":
				errorMessages = append(errorMessages, fmt.Errorf("invalid length %q at %q must be greater or equal than %s",
					e.Value(), location, e.Param()))
			case "resource-name":
				errorMessages = append(errorMessages, fmt.Errorf("invalid resource name %q at %q",
					e.Value(), location))
				hasResNameErr = true
			case "endswith":
				errorMessages = append(errorMessages, fmt.Errorf("invalid value %q at %q must end with %s",
					e.Value(), location, e.Param()))
			case "env":
				errorMessages = append(errorMessages,
					fmt.Errorf("invalid value %q at env var %q must contain only capital letters, numbers, and underscores",
						e.Value(), location))
			case "krt-version":
				errorMessages = append(errorMessages, fmt.Errorf("invalid value %q at krtVersion %q",
					e.Value(), location))
			default:
				errorMessages = append(errorMessages, fmt.Errorf("%s", e))
			}
		}

		if hasResNameErr {
			errorMessages = append(errorMessages,
				errors.New("the resource names must contain only lowercase alphanumeric characters or '-', e.g. my-resource-name"))
		}

		return errorMessages
	}
	return []error{errors.New("internal error parsing fields validator error messages")} //nolint:wsl
}

func (v *CustomValidator) getWorkflowsValidationErrors(workflows []Workflow) []error {
	var validationErrors []error

	for _, workflow := range workflows {
		exitpointError := v.validateExitpoint(workflow)
		if exitpointError != nil {
			validationErrors = append(validationErrors, exitpointError)
		}
	}

	return validationErrors
}

//nolint:goerr113
func (v *CustomValidator) validateExitpoint(workflow Workflow) error {
	if workflow.Exitpoint == "" {
		return fmt.Errorf("missing exitpoint in workflow")
	}

	if !v.isNodeDefined(workflow.Nodes, workflow.Exitpoint) {
		return fmt.Errorf("exitpoint node %q not found in workflow %q nodes", workflow.Exitpoint, workflow.Name)
	}

	return nil
}

func (v *CustomValidator) isNodeDefined(nodes []Node, nodeToFind string) bool {
	for _, node := range nodes {
		if node.Name == nodeToFind {
			return true
		}
	}

	return false
}

//nolint:goerr113
func (v *CustomValidator) ValidateSrcPaths(krtFile *File, dstDir string) error {
	var errs []string

	for _, workflow := range krtFile.Workflows {
		for _, node := range workflow.Nodes {
			nodeFile := path.Join(dstDir, node.Src)
			if !fileExists(nodeFile) {
				errs = append(errs, newSrcPathValidationError(workflow.Name, node.Name, node.Src))
			}
		}
	}

	return errors.New(strings.Join(errs, "\n"))
}

func (v *CustomValidator) ValidateVersionName(versionName string) bool {
	reResourceName := regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	return reResourceName.MatchString(versionName)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
