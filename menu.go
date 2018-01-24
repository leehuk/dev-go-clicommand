package clicommand

import(
    "strings"
)

func (cmd *CLICommand) AddMenu(name string, desc string, f CLICommandFunc) *CLICommand {
    subcmd := New(name, desc)
    subcmd.parent = cmd

    if f != nil {
        subcmd.f = f
    }

    cmd.children = append(cmd.children, subcmd)

    return subcmd
}

func (cmd *CLICommand) GetMenu(name string) *CLICommand {
    for _, v := range cmd.children {
        if strings.EqualFold(v.name, name) {
            return v
        }
    }

    return nil
}

func (cmd *CLICommand) GetMenuNameChain() string {
    name := cmd.name
    if cmd.parent != nil {
        parentname := cmd.parent.GetMenuNameChain()
        if parentname != "" {
            name = parentname + " " + name
        }
    }
    return name
}

func (cmd *CLICommand) GetMenuNameTop() string {
    if cmd.parent != nil {
        return cmd.parent.GetMenuNameTop()
    } else {
        return cmd.name
    }
}


