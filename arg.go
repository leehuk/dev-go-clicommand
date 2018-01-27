package clicommand

import (
)

type Arg struct {
	name     string
	desc     string
	param    bool
	required bool
	parents  []*Command
}

func NewArg(name string, desc string, param bool) *Arg {
	arg := &Arg{
		name:  name,
		desc:  desc,
		param: param,
	}

	return arg
}

func (arg *Arg) BindCommand(cmd *Command) {
	arg.parents = append(arg.parents, cmd)
	cmd.args = append(cmd.args, arg)
}

func (arg *Arg) SetRequired() *Arg {
	arg.required = true
	return arg
}
