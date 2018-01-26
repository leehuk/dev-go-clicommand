package clicommand

func (cmd *Command) BindCallback(handler Handler) {
	cmd.callbacks = append(cmd.callbacks, handler)
}

func (cmd *Command) RunCallbacks(data *Data) error {
	if len(cmd.callbacks) > 0 {
		for _, handler := range cmd.callbacks {
			if error := handler(data); error != nil {
				return error
			}
		}
	}

	if cmd.parent != nil {
		return cmd.parent.RunCallbacks(data)
	}

	return nil
}
