package productconfiguration

func (c *Handler) ListConfiguration(productID, workflowID, processID, scope string) error {
	config, err := c.productConfig.GetConfiguration(productID)
	if err != nil {
		return err
	}

	var productConf map[string]string

	switch scope {
	case "version":
		productConf = config.GetVersionConfiguration()
	case "workflow":
		productConf, err = config.GetWorkflowConfiguration(workflowID)
		if err != nil {
			return err
		}
	case "process":
		productConf, err = config.GetProcessConfiguration(workflowID, processID)
		if err != nil {
			return err
		}
	default:
		return ErrScopeNotValid
	}

	c.renderer.RenderConfiguration(scope, productConf)

	return nil
}
