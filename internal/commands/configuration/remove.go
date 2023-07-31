package productconfiguration

type RemoveConfigurationOpts struct {
	ProductID  string
	WorkflowID string
	ProcessID  string
	Scope      string
	Key        string
}

func (c *Handler) RemoveConfiguration(opts *RemoveConfigurationOpts) error {
	config, err := c.productConfig.GetConfiguration(opts.ProductID)
	if err != nil {
		return err
	}

	var updatedConf map[string]string

	switch opts.Scope {
	case _versionScope:
		updatedConf = config.DeleteVersionConfig(opts.Key)
	case _workflowScope:
		updatedConf, err = config.DeleteWorkflowConfig(opts.WorkflowID, opts.Key)
		if err != nil {
			return err
		}
	case _processScope:
		updatedConf, err = config.DeleteProcessConfig(opts.WorkflowID, opts.ProcessID, opts.Key)
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

	return nil
}
