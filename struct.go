package clicommand

type CLICommand struct {
    name string
    desc string
    f CLICommandFunc
    parent *CLICommand
    children []*CLICommand
    args []*CLICommandArg
    callbacks []CLICommandFunc
}

type CLICommandArg struct {
    name string
    desc string
    param bool
}

type CLICommandFunc func(*CLICommandData) (err error)

type CLICommandData struct {
    Cmd *CLICommand
    Options map[string]string
    Params []string
}

