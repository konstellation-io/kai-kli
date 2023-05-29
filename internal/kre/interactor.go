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

func (i *Interactor) ListVersions(productName string) error {
	list, err := i.client.Version().List(productName)
	if err != nil {
		return fmt.Errorf("error listing versions: %w", err)
	}

	i.renderer.RenderVersions(list)

	return nil
}

func (i *Interactor) CreateVersion(productName, krt string) error {
	versionName, err := i.client.Version().Create(productName, krt)
	if err != nil {
		return fmt.Errorf("error uploading KRT: %w", err)
	}

	i.logger.Success(fmt.Sprintf("Upload KRT completed, version %q created.", versionName))

	return nil
}

func (i *Interactor) StartVersion(productName, versionName, comment string) error {
	err := i.client.Version().Start(productName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Starting version %q on product %q.", versionName, productName))

	return nil
}

func (i *Interactor) StopVersion(productName, versionName, comment string) error {
	err := i.client.Version().Stop(productName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Stopping version %q on product %q.", versionName, productName))

	return nil
}

func (i *Interactor) PublishVersion(productName, versionName, comment string) error {
	err := i.client.Version().Publish(productName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Publishing version %q on product %q.", versionName, productName))

	return nil
}

func (i *Interactor) UnpublishVersion(productName, versionName, comment string) error {
	err := i.client.Version().Unpublish(productName, versionName, comment)
	if err != nil {
		return err
	}

	i.logger.Success(fmt.Sprintf("Unpublishing version %q on product %q.", versionName, productName))

	return nil
}

func (i *Interactor) ListVersionConfig(productName, versionName string, showValues bool) error {
	cfg, err := i.client.Version().GetConfig(productName, versionName)
	if err != nil {
		return err
	}

	if len(cfg.Vars) == 0 {
		i.logger.Info(fmt.Sprintf("No config found for version %q on product %q.", versionName, productName))
		return nil
	}

	i.renderer.RenderVars(cfg, showValues)

	return nil
}

func (i *Interactor) UpdateVersionConfig(productName, versionName string, newConfig []version.ConfigVariableInput) error {
	completed, err := i.client.Version().UpdateConfig(productName, versionName, newConfig)
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

func (i *Interactor) ListProducts() error {
	products, err := i.client.Product().List()
	if err != nil {
		return fmt.Errorf("error getting products: %w", err)
	}

	i.renderer.RenderProducts(products)

	return nil
}

func (i *Interactor) CreateProduct(name, description string) error {
	err := i.client.Product().Create(name, description)
	if err != nil {
		return fmt.Errorf("error creating product: %w", err)
	}

	return nil
}
