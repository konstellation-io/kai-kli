package productconfiguration

import (
	productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"
)

func (c *Handler) AddConfiguration(productID, workflowID, processID, scope, key, value string) error {
	config, err := c.productConfig.GetConfiguration(productID)
	if err != nil {
		return err
	}

	var updatedConf map[string]string

	switch scope {
	case "version":
		updatedConf = config.UpdateVersionConfig(productconfiguration.ConfigProperty{Key: key, Value: value})
	case "workflow":
		updatedConf, err = config.UpdateWorkflowConfig(workflowID, productconfiguration.ConfigProperty{Key: key, Value: value})
		if err != nil {
			return err
		}
	case "process":
		updatedConf, err = config.UpdateProcessConfig(workflowID, processID, productconfiguration.ConfigProperty{Key: key, Value: value})
		if err != nil {
			return err
		}
	default:
		return ErrScopeNotValid
	}

	err = c.productConfig.WriteConfiguration(config, productID)
	if err != nil {
		return err
	}

	c.renderer.RenderConfiguration(scope, updatedConf)

	return nil
}
