package clicommand

import (
	"fmt"
	"strings"
)

func NewArg(name string, desc string, param bool) *Arg {
	arg := &Arg{
		name:  name,
		desc:  desc,
		param: param,
	}

	return arg
}

func (cmd *Command) BindArg(arg ...*Arg) {
	cmd.args = append(cmd.args, arg...)
}

func (cmd *Command) NewArg(name string, desc string, param bool) *Arg {
	arg := NewArg(name, desc, param)
	cmd.BindArg(arg)
	return arg
}

func (arg *Arg) SetRequired() *Arg {
	arg.required = true
	return arg
}

func (cmd *Command) GetArg(name string, param bool) *Arg {
	for _, v := range cmd.args {
		if strings.EqualFold(v.name, name) && v.param == param {
			return v
		}
	}

	// not found, may be a parameter to a parent menu
	if cmd.parent != nil {
		return cmd.parent.GetArg(name, param)
	}

	return nil
}

func (cmd *Command) HasRequiredArgs(data *Data) error {
	for _, v := range cmd.args {
		if _, ok := data.Options[v.name]; v.required && !ok {
			return fmt.Errorf("Required option missing: %s", v.name)

		}
	}

	if cmd.parent != nil {
		return cmd.parent.HasRequiredArgs(data)
	}

	return nil
}
