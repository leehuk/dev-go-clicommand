package clicommand

import (
	"strings"
)

func (cmd *Command) AddCommand(name string, desc string, handler Handler) *Command {
	subcmd := New(name, desc)
	subcmd.parent = cmd

	if handler != nil {
		subcmd.handler = handler
	}

	cmd.children = append(cmd.children, subcmd)

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
