package clicommand

import ()

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

func (arg *Arg) UnbindCommand(cmd *Command) {
	var newparents []*Command
	var newargs []*Arg

	for _, cmdv := range arg.parents {
		if cmdv != cmd {
			newparents = append(newparents, cmdv)
		}
	}

	for _, argv := range cmd.args {
		if argv != arg {
			newargs = append(newargs, argv)
		}
	}

	arg.parents = newparents
	cmd.args = newargs
}

func (arg *Arg) GetRequired() bool {
	return arg.required
}

func (arg *Arg) SetRequired() *Arg {
	arg.required = true
	return arg
}

func (arg *Arg) GetParents() []*Command {
	return arg.parents
}
