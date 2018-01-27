package clicommand

import ()

type Option struct {
	name     string
	desc     string
	param    bool
	required bool
	parents  []*Command
}

func NewOption(name string, desc string, param bool) *Option {
	opt := &Option{
		name:  name,
		desc:  desc,
		param: param,
	}

	return opt
}

func (self *Option) BindCommand(cmd *Command) {
	self.parents = append(self.parents, cmd)
	cmd.options = append(cmd.options, self)
}

func (self *Option) UnbindCommand(cmd *Command) {
	var newparents []*Command
	var newoptions []*Option

	for _, cmdi := range self.parents {
		if cmdi != cmd {
			newparents = append(newparents, cmdi)
		}
	}

	for _, opt := range cmd.options {
		if opt != self {
			newoptions = append(newoptions, opt)
		}
	}

	self.parents = newparents
	cmd.options = newoptions
}

func (self *Option) GetRequired() bool {
	return self.required
}

func (self *Option) SetRequired() *Option {
	self.required = true
	return self
}

func (self *Option) GetParents() []*Command {
	return self.parents
}
