package clicommand

func (cmd *CLICommand) AddCallback(f CLICommandFunc) {
    cmd.callbacks = append(cmd.callbacks, f)
}

func (cmd *CLICommand) RunCallbacks(data *CLICommandData) error {
    if len(cmd.callbacks) > 0 {
        for _, f := range cmd.callbacks {
            if error := f(data); error != nil {
                return error
            }
        }
    }

    if cmd.parent != nil {
        return cmd.parent.RunCallbacks(data)
    }

    return nil
}
