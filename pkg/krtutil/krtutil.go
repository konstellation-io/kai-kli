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
	krtProcesses := make([]krt.Process, 0, len(processes))

	for _, process := range processes {
		replicaCopy := process.Replicas
		gpuCopy := process.GPU

		krtProcesses = append(krtProcesses, krt.Process{
			Name:           process.Name,
			Type:           krt.ProcessType(process.Type),
			Image:          process.Image,
			Replicas:       &replicaCopy,
			GPU:            &gpuCopy,
			Config:         mapConfigToKrt(process.Config),
			ObjectStore:    mapObjectStoreToKrt(process.ObjectStore),
			Secrets:        mapConfigToKrt(process.Secrets),
			Subscriptions:  process.Subscriptions,
			Networking:     mapNetworkingToKrt(process.Networking),
			ResourceLimits: mapResourceLimitsToKrt(process.ResourceLimits),
			NodeSelectors:  process.NodeSelectors,
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
