//go:build unit

package process_test

import (
	"fmt"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/krt/pkg/krt"

	configservice "github.com/konstellation-io/kli/internal/services/product_configuration"
)

type processMatcher struct {
	processes []krt.Process
}

func ProcessMatcher(values ...krt.Process) gomock.Matcher {
	return &processMatcher{processes: values}
}

func (pm *processMatcher) Matches(arg interface{}) bool {
	actual, ok := arg.([]krt.Process)
	if !ok {
		return false
	}

	for i, p := range pm.processes { //nolint:gocritic
		return p.Name == actual[i].Name &&
			p.Type == actual[i].Type &&
			p.Image == actual[i].Image &&
			*p.Replicas == *actual[i].Replicas &&
			*p.GPU == *actual[i].GPU &&
			reflect.DeepEqual(p.Config, actual[i].Config) &&
			reflect.DeepEqual(p.ObjectStore, actual[i].ObjectStore) &&
			reflect.DeepEqual(p.Secrets, actual[i].Secrets) &&
			reflect.DeepEqual(p.Subscriptions, actual[i].Subscriptions) &&
			reflect.DeepEqual(p.Subscriptions, actual[i].Subscriptions)
	}

	return false
}

func (pm *processMatcher) String() string {
	return fmt.Sprintf("%v", pm.processes)
}

func _getDefaultKaiConfig() *configservice.KaiProductConfiguration {
	return &configservice.KaiProductConfiguration{
		Krt: &krt.Krt{
			Version:   "v0.0.1",
			Config:    map[string]string{"test1": "value1"},
			Workflows: []krt.Workflow{_getDefaultWorkflow()},
		},
	}
}

func _getDefaultWorkflow() krt.Workflow {
	return krt.Workflow{
		Name:      "Workflow1",
		Type:      "data",
		Config:    map[string]string{"test2": "value2"},
		Processes: []krt.Process{_getDefaultProcess()},
	}
}

func _getDefaultProcess() krt.Process {
	replicas := 1
	gpu := false

	return krt.Process{
		Name:          "process-1",
		Type:          krt.ProcessTypeTrigger,
		Image:         "kst/trigger",
		Replicas:      &replicas,
		GPU:           &gpu,
		Config:        map[string]string{"test3": "value3"},
		ObjectStore:   nil,
		Secrets:       map[string]string{},
		Subscriptions: []string{"subject1", "subject2"},
		ResourceLimits: &krt.ProcessResourceLimits{
			CPU: &krt.ResourceLimit{
				Request: "500m",
				Limit:   "1.5",
			},
			Memory: &krt.ResourceLimit{
				Request: "64Mi",
				Limit:   "128Mi",
			},
		},
		Networking: &krt.ProcessNetworking{
			TargetPort:      20000,
			DestinationPort: 21000,
			Protocol:        krt.NetworkingProtocolHTTP,
		},
	}
}
