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
    cmd *CLICommand
    options map[string]string
    params []string
}

