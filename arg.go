package clicommand

import(
    "strings"
)

func (cmd *Command) AddArg(name string, desc string, param bool)  *Arg {
    arg := &Arg{
        name,
        desc,
        param,
    }
    cmd.args = append(cmd.args, arg)
    return arg
}

func (cmd *Command) GetArg(name string, param bool) *Arg {
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


