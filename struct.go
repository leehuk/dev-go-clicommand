package clicommand

type Command struct {
    name string
    desc string
    f CommandFunc
    parent *Command
    children []*Command
    args []*CommandArg
}

type CommandArg struct {
    name string
    desc string
    param bool
}

type CommandFunc func([]string, map[string]string) (err error)
