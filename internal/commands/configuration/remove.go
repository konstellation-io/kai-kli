package productconfiguration

func (c *Handler) RemoveConfiguration(productID, workflowID, processID, scope, key string) error {
	config, err := c.productConfig.GetConfiguration(productID)
	if err != nil {
		return err
	}

	var updatedConf map[string]string

	switch scope {
	case "version":
		updatedConf = config.DeleteVersionConfig(key)
	case "workflow":
		updatedConf, err = config.DeleteWorkflowConfig(workflowID, key)
		if err != nil {
			return err
		}
	case "process":
		updatedConf, err = config.DeleteProcessConfig(workflowID, processID, key)
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
