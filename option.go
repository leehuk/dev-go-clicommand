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

func (self *Arg) BindCommand(cmd *Command) {
	self.parents = append(self.parents, cmd)
	cmd.args = append(cmd.args, self)
}

func (self *Arg) UnbindCommand(cmd *Command) {
	var newparents []*Command
	var newargs []*Arg

	for _, cmdi := range self.parents {
		if cmdi != cmd {
			newparents = append(newparents, cmdi)
		}
	}

	for _, arg := range cmd.args {
		if arg != self {
			newargs = append(newargs, arg)
		}
	}

	self.parents = newparents
	cmd.args = newargs
}

func (self *Arg) GetRequired() bool {
	return self.required
}

func (self *Arg) SetRequired() *Arg {
	self.required = true
	return self
}

func (self *Arg) GetParents() []*Command {
	return self.parents
}
