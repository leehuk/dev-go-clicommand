package clicommand

func (cmd *Command) BindPreCallback(handler Handler) {
	cmd.precallbacks = append(cmd.precallbacks, handler)
}

func (cmd *Command) BindValidateCallback(handler Handler) {
	cmd.valcallbacks = append(cmd.valcallbacks, handler)
}

func (cmd *Command) RunPreCallbacks(data *Data) error {
	if len(cmd.precallbacks) > 0 {
		for _, handler := range cmd.precallbacks {
			if error := handler(data); error != nil {
				return error
			}
		}
	}

	if cmd.parent != nil {
		return cmd.parent.RunPreCallbacks(data)
	}

	return nil
}

func (cmd *Command) RunValidateCallbacks(data *Data) error {
	if len(cmd.valcallbacks) > 0 {
		for _, handler := range cmd.valcallbacks {
			if error := handler(data); error != nil {
				return error
			}
		}
	}

	if cmd.parent != nil {
		return cmd.parent.RunValidateCallbacks(data)
	}

	return nil
}
