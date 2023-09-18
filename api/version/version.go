package version

import (
	"time"
)

type Version struct {
	ID          string                  `json:"id"`
	Tag         string                  `json:"tag"`
	Description string                  `json:"description"`
	Config      []ConfigurationVariable `json:"config"`
	Workflows   []Workflow              `json:"workflows"`

	CreationDate   time.Time `json:"creationDate"`
	CreationAuthor string    `json:"creationAuthor"`

	PublicationDate   *time.Time `json:"publicationDate"`
	PublicationAuthor *string    `json:"publicationAuthor"`

	Status Status `json:"status"`
	Error  string `json:"error"`
}

type ConfigurationVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Status string

const (
	VersionStatusCreating  Status = "CREATING"
	VersionStatusCreated   Status = "CREATED"
	VersionStatusStarting  Status = "STARTING"
	VersionStatusStarted   Status = "STARTED"
	VersionStatusPublished Status = "PUBLISHED"
	VersionStatusStopping  Status = "STOPPING"
	VersionStatusStopped   Status = "STOPPED"
	VersionStatusError     Status = "ERROR"
)

type Workflow struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	Type      WorkflowType            `json:"type"`
	Config    []ConfigurationVariable `json:"config"`
	Processes []Process               `json:"processes"`
}

type WorkflowType string

const (
	WorkflowTypeData     WorkflowType = "data"
	WorkflowTypeTraining WorkflowType = "training"
	WorkflowTypeFeedback WorkflowType = "feedback"
	WorkflowTypeServing  WorkflowType = "serving"
)

type Process struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Type          ProcessType             `json:"type"`
	Image         string                  `json:"image"`
	Replicas      int32                   `json:"replicas"`
	GPU           bool                    `json:"gpu"`
	Config        []ConfigurationVariable `json:"config"`
	Secrets       []ConfigurationVariable `json:"secrets"`
	Subscriptions []string                `json:"subscriptions"`
	Status        ProcessStatus           `json:"status"`
}

type ProcessType string

const (
	ProcessTypeTrigger ProcessType = "trigger"
	ProcessTypeTask    ProcessType = "task"
	ProcessTypeExit    ProcessType = "exit"
)

type ProcessStatus string

const (
	ProcessStatusStarting ProcessStatus = "STARTING"
	ProcessStatusStarted  ProcessStatus = "STARTED"
	ProcessStatusStopped  ProcessStatus = "STOPPED"
	ProcessStatusError    ProcessStatus = "ERROR"
)
