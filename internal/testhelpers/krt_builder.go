package testhelpers

import (
	"gopkg.in/yaml.v3"

	"github.com/konstellation-io/kli/internal/commands/krt/entity"
)

type KrtBuilder struct {
	krtYaml *entity.File
}

func NewKrtBuilder() *KrtBuilder {
	return &KrtBuilder{
		krtYaml: &entity.File{
			KrtVersion:  "v2",
			Version:     "version-name",
			Description: "Test description",
			Entrypoint: &entity.Entrypoint{
				Proto: "valid.proto",
				Image: "test/image",
			},
			Config: &entity.Config{
				Variables: []string{"TEST_VAR"},
				Files:     []string{"TEST_FILE"},
			},
			Workflows: []entity.Workflow{
				{
					Name:       "valid-workflow",
					Entrypoint: "valid-entrypoint",
					Exitpoint:  "test-node",
					Nodes: []entity.Node{
						{
							Name:          "test-node",
							Image:         "test/image",
							Src:           "src/test",
							GPU:           false,
							Replicas:      1,
							Subscriptions: []string{"entrypoint"},
						},
					},
				},
			},
		},
	}
}

func (k *KrtBuilder) WithKrtVersion(krtVersion string) *KrtBuilder {
	k.krtYaml.KrtVersion = krtVersion
	return k
}

func (k *KrtBuilder) WithVersion(version string) *KrtBuilder {
	v, _ := entity.NewResourceName(version)
	k.krtYaml.Version = v

	return k
}

func (k *KrtBuilder) WithDescription(description string) *KrtBuilder {
	k.krtYaml.Description = description
	return k
}

func (k *KrtBuilder) WithEntrypoint(entrypoint *entity.Entrypoint) *KrtBuilder {
	k.krtYaml.Entrypoint = entrypoint
	return k
}

func (k *KrtBuilder) WithEntrypointProto(proto string) *KrtBuilder {
	k.krtYaml.Entrypoint.Proto = proto
	return k
}

func (k *KrtBuilder) WithEntrypointImage(image string) *KrtBuilder {
	k.krtYaml.Entrypoint.Image = image
	return k
}

func (k *KrtBuilder) WithConfigVars(vars []string) *KrtBuilder {
	k.krtYaml.Config.Variables = vars
	return k
}

func (k *KrtBuilder) WithWorkflows(workflows []entity.Workflow) *KrtBuilder {
	k.krtYaml.Workflows = workflows
	return k
}

func (k *KrtBuilder) Build() *entity.File {
	return k.krtYaml
}

func (k *KrtBuilder) AsString() string {
	krtFile, err := yaml.Marshal(k.krtYaml)
	if err != nil {
		return ""
	}

	return string(krtFile)
}
