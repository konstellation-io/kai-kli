package kre

import (
	"fmt"
	"io"

	api "github.com/konstellation-io/kli/api/kre"
	"github.com/konstellation-io/kli/api/kre/version"
	"github.com/konstellation-io/kli/internal/logging"
	"github.com/konstellation-io/kli/internal/render"
)

type Interactor struct {
	logger   logging.Interface
	client   api.KreInterface
	renderer render.Renderer
}

func NewInteractor(
	logger logging.Interface,
	client api.KreInterface,
	renderer render.Renderer,
) *Interactor {
	return &Interactor{
		logger:   logger,
		client:   client,
		renderer: renderer,
	}
}

func NewInteractorWithDefaultRenderer(
	logger logging.Interface,
	client api.KreInterface,
	writer io.Writer,
) *Interactor {
	renderer := render.NewDefaultCliRenderer(logger, writer)
	return NewInteractor(logger, client, renderer)
}

func (i *Interactor) ListVersions(runtimeName string) error {
	list, err := i.client.Version().List(runtimeName)
	if err != nil {
		return fmt.Errorf("error listing versions: %w", err)
	}

	i.renderer.RenderVersions(list)

	return nil
}

func (i *Interactor) CreateVersion(runtimeName, krt string) error {
	versionName, err := i.client.Version().Create(runtimeName, krt)
	if err != nil {
		return fmt.Errorf("error uploading KRT: %w", err)
	}

	i.logger.Success(fmt.Sprintf("Upload KRT completed, version %q created.", versionName))

	return nil
}

func (i *Interactor) StartVersion(runtimeName, versionName, comment string) error {
	err := i.client.Version().Start(runtimeName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Starting version %q on runtime %q.", versionName, runtimeName))

	return nil
}

func (i *Interactor) StopVersion(runtimeName, versionName, comment string) error {
	err := i.client.Version().Stop(runtimeName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Stopping version %q on runtime %q.", versionName, runtimeName))

	return nil
}

func (i *Interactor) PublishVersion(runtimeName, versionName, comment string) error {
	err := i.client.Version().Publish(runtimeName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Publishing version %q on runtime %q.", versionName, runtimeName))

	return nil
}

func (i *Interactor) UnpublishVersion(runtimeName, versionName, comment string) error {
	err := i.client.Version().Unpublish(runtimeName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Unpublishing version %q on runtime %q.", versionName, runtimeName))

	return nil
}

func (i *Interactor) ListVersionConfig(runtimeName, versionName string, showValues bool) error {
	cfg, err := i.client.Version().GetConfig(runtimeName, versionName)
	if err != nil {
		return err
	}

	if len(cfg.Vars) == 0 {
		i.logger.Info(fmt.Sprintf("No config found for version %q on runtime %q.", versionName, runtimeName))
		return nil
	}

	i.renderer.RenderVars(cfg, showValues)

	return nil
}

func (i *Interactor) UpdateVersionConfig(runtimeName, versionName string, newConfig []version.ConfigVariableInput) error {
	completed, err := i.client.Version().UpdateConfig(runtimeName, versionName, newConfig)
	if err != nil {
		return err
	}

	status := "updated"
	if completed {
		status = "completed"
	}

	i.logger.Success(fmt.Sprintf("Config %s for version %q.", status, versionName))

	return nil
}

func (i *Interactor) ListRuntimes() error {
	runtimes, err := i.client.Runtime().List()
	if err != nil {
		return fmt.Errorf("error getting runtimes: %w", err)
	}

	i.renderer.RenderRuntimes(runtimes)

	return nil
}

func (i *Interactor) CreateRuntime(name, description string) error {
	err := i.client.Runtime().Create(name, description)
	if err != nil {
		return fmt.Errorf("error creating runtime: %w", err)
	}

	return nil
}
