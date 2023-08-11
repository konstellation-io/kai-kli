package productconfiguration

type ListConfigurationOpts struct {
	ProductID  string
	WorkflowID string
	ProcessID  string
	Scope      string
}

func (c *Handler) ListConfiguration(opts *ListConfigurationOpts) error {
	config, err := c.productConfig.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	var productConf map[string]string

	switch opts.Scope {
	case _versionScope:
		productConf = config.GetVersionConfiguration()
	case _workflowScope:
		productConf, err = config.GetWorkflowConfiguration(opts.WorkflowID)
		if err != nil {
			return err
		}
	case _processScope:
		productConf, err = config.GetProcessConfiguration(opts.WorkflowID, opts.ProcessID)
		if err != nil {
			return err
		}
	default:
		return ErrScopeNotValid
	}

	c.renderer.RenderConfiguration(opts.Scope, productConf)

	return nil
}
