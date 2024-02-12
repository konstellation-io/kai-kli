package entity

import "time"

const (
	VersionStatusStarted   = "STARTED"
	VersionStatusPublished = "PUBLISHED"
	VersionStatusStopped   = "STOPPED"
	VersionStatusCritical  = "CRITICAL"
)

type ConfigurationVariable struct {
	Key   string
	Value string
}

type Version struct {
	Tag               string                  `json:"tag"`
	Description       string                  `json:"description"`
	Config            []ConfigurationVariable `json:"config"`
	Workflows         []Workflow              `json:"workflows"`
	PublishedTriggers []TriggerEndpoint       `json:"publishedTriggers"`

	CreationDate time.Time `json:"creationDate"`
	Status       string    `json:"status"`

	Error string `json:"error"`
}

type Workflow struct {
	Name      string                  `json:"name"`
	Type      string                  `json:"type"`
	Config    []ConfigurationVariable `json:"config"`
	Processes []Process               `json:"processes"`
}

type Process struct {
	Name           string                  `json:"name"`
	Type           string                  `json:"type"`
	Image          string                  `json:"image"`
	Replicas       int                     `json:"replicas"`
	GPU            bool                    `json:"gpu"`
	Config         []ConfigurationVariable `json:"config"`
	ObjectStore    *ProcessObjectStore     `json:"objectStore"`
	Secrets        []ConfigurationVariable `json:"secrets"`
	Subscriptions  []string                `json:"subscriptions"`
	Networking     *ProcessNetworking      `json:"networking"`
	ResourceLimits *ProcessResourceLimits  `json:"resourceLimits"`
	NodeSelectors  map[string]string       `json:"nodeSelectors"`
}

type ProcessObjectStore struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}

type ProcessNetworking struct {
	TargetPort      int    `json:"targetPort"`
	DestinationPort int    `json:"destinationPort"`
	Protocol        string `json:"protocol"`
}

type ResourceLimit struct {
	Request string `json:"request"`
	Limit   string `json:"limit"`
}

type ProcessResourceLimits struct {
	CPU    *ResourceLimit `json:"cpu"`
	Memory *ResourceLimit `json:"memory"`
}
