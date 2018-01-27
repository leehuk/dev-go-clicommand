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

func (o *Option) BindCommand(cmd *Command) {
	o.parents = append(o.parents, cmd)
	cmd.options = append(cmd.options, o)
}

func (o *Option) UnbindCommand(cmd *Command) {
	var newparents []*Command
	var newoptions []*Option

	for _, cmdi := range o.parents {
		if cmdi != cmd {
			newparents = append(newparents, cmdi)
		}
	}

	for _, opt := range cmd.options {
		if opt != o {
			newoptions = append(newoptions, opt)
		}
	}

	o.parents = newparents
	cmd.options = newoptions
}

func (o *Option) GetRequired() bool {
	return o.required
}

func (o *Option) SetRequired() *Option {
	o.required = true
	return o
}

func (o *Option) GetParents() []*Command {
	return o.parents
}
