package clicommand

import (
	"strings"
)

func NewCommand(name string, desc string, handler Handler) *Command {
	cmd := &Command{
		name:    name,
		desc:    desc,
		handler: handler,
	}

	return cmd
}

func (cmd *Command) BindCommand(subcmd ...*Command) {
	cmd.children = append(cmd.children, subcmd...)
	for _, v := range subcmd {
		v.parent = cmd
	}
}

func (cmd *Command) NewCommand(name string, desc string, handler Handler) *Command {
	subcmd := NewCommand(name, desc, handler)
	cmd.BindCommand(subcmd)
	return subcmd
}

func (cmd *Command) GetCommand(name string) *Command {
	for _, v := range cmd.children {
		if strings.EqualFold(v.name, name) {
			return v
		}
	}

	return nil
}

func (cmd *Command) GetCommandNameChain() string {
	name := cmd.name
	if cmd.parent != nil {
		parentname := cmd.parent.GetCommandNameChain()
		if parentname != "" {
			name = parentname + " " + name
		}
	}
	return name
}

func (cmd *Command) GetCommandNameTop() string {
	if cmd.parent != nil {
		return cmd.parent.GetCommandNameTop()
	}

	return cmd.name
}
