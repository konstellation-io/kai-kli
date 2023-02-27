package entity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/konstellation-io/kli/internal/krt/errors"
	"github.com/konstellation-io/kli/pkg/collectionutils"
)

const resourceNameErrorMessage = "Invalid resource name; only numbers, hyphens and lowercase letters are allowed"

func (f *File) Validate() []*errors.ValidationError {
	var validationErrors []*errors.ValidationError

	err := f.ValidateKrtVersion()
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	err = f.ValidateVersion()
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	err = f.ValidateDescription()
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	err = f.ValidateConfig()
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	err = f.ValidateEntrypoint()
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	errs := f.ValidateWorkflows()
	if errs != nil {
		validationErrors = append(validationErrors, errs...)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (f *File) ValidateKrtVersion() *errors.ValidationError {
	validVersions := []string{"v2"}
	if !collectionutils.ArrayContains(validVersions, f.KrtVersion) {
		return errors.NewValidationError(
			"krtVersion",
			fmt.Sprintf("Invalid KRT Version, valid versions are: %q", validVersions),
		)
	}

	return nil
}

func (f *File) ValidateVersion() *errors.ValidationError {
	if !IsValidResourceName(string(f.Version)) {
		return errors.NewValidationError(
			"version",
			resourceNameErrorMessage,
		)
	}

	return nil
}

func (f *File) ValidateDescription() *errors.ValidationError {
	if f.Description == "" {
		return errors.NewValidationError(
			"description",
			"A description must be given",
		)
	}

	return nil
}

func (f *File) ValidateEntrypoint() *errors.ValidationError {
	isValidEntrypoint := f.Entrypoint != nil && f.Entrypoint.Image != "" && f.Entrypoint.Proto != ""
	if !isValidEntrypoint {
		return errors.NewValidationError(
			"entrypoint",
			"An entrypoint definition [Image, Proto] must be given",
		)
	}

	return nil
}

func (f *File) ValidateConfig() *errors.ValidationError {
	if f.Config == nil {
		return nil
	}

	reEnvVar := regexp.MustCompile("^[A-Z0-9]([_A-Z0-9]*[A-Z0-9])?$")

	var invalidVars []string

	for i, envVar := range f.Config.Variables {
		if !reEnvVar.MatchString(envVar) {
			invalidVars = append(invalidVars, strconv.Itoa(i))
		}
	}

	if len(invalidVars) != 0 {
		return errors.NewValidationError(
			"config.variables",
			fmt.Sprintf(
				"Following fields don't match the expected format (Uppercase letters, numbers and low bars): %s",
				strings.Join(invalidVars, ", "),
			),
		)
	}

	return nil
}

func (f *File) ValidateWorkflows() []*errors.ValidationError {
	if len(f.Workflows) == 0 {
		return []*errors.ValidationError{
			{
				Field:   "workflows",
				Message: "At least 1 workflow must be defined",
			},
		}
	}

	var validationErrors []*errors.ValidationError

	workflowsNames := make([]string, 0, len(f.Workflows))

	for _, workflow := range f.Workflows {
		workflowsNames = append(workflowsNames, string(workflow.Name))
		err := f.ValidateWorkflow(workflow)
		validationErrors = append(validationErrors, err...)
	}

	namesSet := collectionutils.ToSet(workflowsNames)
	if len(namesSet) != len(f.Workflows) {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"worklfow.name",
			"Workflows names must be unique",
		))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (f *File) ValidateWorkflow(workflow Workflow) []*errors.ValidationError {
	var validationErrors []*errors.ValidationError
	if !IsValidResourceName(string(workflow.Name)) {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"workflow.name",
			resourceNameErrorMessage,
		))
	}

	if workflow.Entrypoint == "" {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"workflow.entrypoint",
			"Entrypoint must be defined",
		))
	}

	isValidExitpoint := false

	for _, node := range workflow.Nodes {
		if string(node.Name) == workflow.Exitpoint {
			isValidExitpoint = true
		}

		errs := f.ValidateNode(node)
		validationErrors = append(validationErrors, errs...)
	}

	if !isValidExitpoint {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"workflow.exitpoint",
			"Exitpoint not defined or not found in nodes list",
		))
	}

	return validationErrors
}

func (f *File) ValidateNode(node Node) []*errors.ValidationError { //nolint:gocritic
	var validationErrors []*errors.ValidationError
	if !IsValidResourceName(string(node.Name)) {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"node.name",
			resourceNameErrorMessage,
		))
	}

	if node.Image == "" {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"node.image",
			"Image must be defined",
		))
	}

	if node.Src == "" {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"node.src",
			"Src must be defined",
		))
	}

	if len(node.Subscriptions) < 1 {
		validationErrors = append(validationErrors, errors.NewValidationError(
			"node.subscriptions",
			"Nodes should have at least one subscription",
		))
	}

	return validationErrors
}

func IsValidResourceName(name string) bool {
	reResourceName := regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	return reResourceName.MatchString(name) && len(name) < 20
}
