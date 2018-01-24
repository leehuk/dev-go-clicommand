package clicommand

import(
    "strings"
)

func (cmd *CLICommand) AddArg(name string, desc string, param bool) {
    arg := &CLICommandArg{
        name,
        desc,
        param,
    }
    cmd.args = append(cmd.args, arg)
}

func (cmd *CLICommand) GetArg(name string, param bool) *CLICommandArg {
    for _, v := range cmd.args {
        if strings.EqualFold(v.name, name) && v.param == param {
            return v
        }
    }

    // not found, may be a parameter to a parent menu
    if cmd.parent != nil {
        return cmd.parent.GetArg(name, param)
    } else {
        return nil
    }
}


