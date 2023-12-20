package krtutil

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/krt/pkg/krt"
)

func MapVersionToKrt(version *entity.Version) *krt.Krt {
	return &krt.Krt{
		Version:     version.Tag,
		Description: version.Description,
		Config:      mapConfigToKrt(version.Config),
		Workflows:   mapWorkflowsToKrt(version.Workflows),
	}
}

func mapConfigToKrt(config []entity.ConfigurationVariable) map[string]string {
	krtConfig := map[string]string{}

	for _, c := range config {
		krtConfig[c.Key] = c.Value
	}

	return krtConfig
}

func mapWorkflowsToKrt(workflows []entity.Workflow) []krt.Workflow {
	krtWorkflows := []krt.Workflow{}

	for _, w := range workflows {
		krtWorkflows = append(krtWorkflows, krt.Workflow{
			Name:      w.Name,
			Type:      krt.WorkflowType(w.Type),
			Config:    mapConfigToKrt(w.Config),
			Processes: mapProcessesToKrt(w.Processes),
		})
	}

	return krtWorkflows
}

func mapProcessesToKrt(processes []entity.Process) []krt.Process {
	krtProcesses := []krt.Process{}

	for _, p := range processes {
		replicaCopy := p.Replicas
		gpuCopy := p.GPU

		krtProcesses = append(krtProcesses, krt.Process{
			Name:           p.Name,
			Type:           krt.ProcessType(p.Type),
			Image:          p.Image,
			Replicas:       &replicaCopy,
			GPU:            &gpuCopy,
			Config:         mapConfigToKrt(p.Config),
			ObjectStore:    mapObjectStoreToKrt(p.ObjectStore),
			Secrets:        mapConfigToKrt(p.Secrets),
			Subscriptions:  p.Subscriptions,
			Networking:     mapNetworkingToKrt(p.Networking),
			ResourceLimits: mapResourceLimitsToKrt(p.ResourceLimits),
		})
	}

	return krtProcesses
}

func mapObjectStoreToKrt(objectStore *entity.ProcessObjectStore) *krt.ProcessObjectStore {
	if objectStore == nil {
		return nil
	}

	return &krt.ProcessObjectStore{
		Name:  objectStore.Name,
		Scope: krt.ObjectStoreScope(objectStore.Scope),
	}
}

func mapNetworkingToKrt(networking *entity.ProcessNetworking) *krt.ProcessNetworking {
	if networking == nil {
		return nil
	}

	return &krt.ProcessNetworking{
		TargetPort:      networking.TargetPort,
		DestinationPort: networking.DestinationPort,
		Protocol:        krt.NetworkingProtocol(networking.Protocol),
	}
}

func mapResourceLimitsToKrt(resourceLimits *entity.ProcessResourceLimits) *krt.ProcessResourceLimits {
	if resourceLimits == nil {
		return nil
	}

	return &krt.ProcessResourceLimits{
		CPU:    mapResourceLimitToKrt(resourceLimits.CPU),
		Memory: mapResourceLimitToKrt(resourceLimits.Memory),
	}
}

func mapResourceLimitToKrt(resourceLimit *entity.ResourceLimit) *krt.ResourceLimit {
	if resourceLimit == nil {
		return nil
	}

	return &krt.ResourceLimit{
		Request: resourceLimit.Request,
		Limit:   resourceLimit.Limit,
	}
}
