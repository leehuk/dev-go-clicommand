package clicommand

type CLICommand struct {
    name string
    desc string
    f CLICommandFunc
    parent *CLICommand
    children []*CLICommand
    args []*CLICommandArg
}

type CLICommandArg struct {
    name string
    desc string
    param bool
}

type CLICommandFunc func([]string, map[string]string) (err error)
