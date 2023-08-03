package productconfiguration

import (
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

type AddConfigurationOpts struct {
	ProductID  string
	WorkflowID string
	ProcessID  string
	Scope      string
	Key        string
	Value      string
}

func (c *Handler) AddConfiguration(opts *AddConfigurationOpts) error {
	config, err := c.productConfig.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	var updatedConf map[string]string

	switch opts.Scope {
	case _versionScope:
		updatedConf = config.UpdateVersionConfig(productconfiguration.ConfigProperty{Key: opts.Key, Value: opts.Value})
	case _workflowScope:
		updatedConf, err = config.UpdateWorkflowConfig(opts.WorkflowID, productconfiguration.ConfigProperty{Key: opts.Key, Value: opts.Value})
		if err != nil {
			return err
		}
	case _processScope:
		updatedConf, err = config.UpdateProcessConfig(opts.WorkflowID, opts.ProcessID,
			productconfiguration.ConfigProperty{Key: opts.Key, Value: opts.Value})
		if err != nil {
			return err
		}
	default:
		return ErrScopeNotValid
	}

	err = c.productConfig.WriteConfiguration(config, opts.ProductID)
	if err != nil {
		return err
	}

	c.renderer.RenderConfiguration(opts.Scope, updatedConf)

	c.logger.Success("Configuration added successfully")

	return nil
}
