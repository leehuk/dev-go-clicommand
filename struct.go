package clicommand

type CLICommandMenu struct {
    name string
    desc string
    f CLICommandFunc
    parent *CLICommandMenu
    children []*CLICommandMenu
    args []*CLICommandArg
}

type CLICommandArg struct {
    name string
    desc string
    param bool
}

type CLICommandFunc func([]string, map[string]string) (err error)
