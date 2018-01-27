package clicommand

import (
	"fmt"
	"strings"
)

type Command struct {
	name         string
	desc         string
	handler      Handler
	parent       *Command
	children     []*Command
	args         []*Arg
	callbackspre []Handler
	callbacks    []Handler
}

func NewCommand(name string, desc string, handler Handler) *Command {
	cmd := &Command{
		name:    name,
		desc:    desc,
		handler: handler,
	}

	return cmd
}

func (cmd *Command) NewCommand(name string, desc string, handler Handler) *Command {
	subcmd := NewCommand(name, desc, handler)
	cmd.BindCommand(subcmd)
	return subcmd
}

func (cmd *Command) BindCommand(subcmd ...*Command) {
	cmd.children = append(cmd.children, subcmd...)
	for _, v := range subcmd {
		v.parent = cmd
	}
}

func (cmd *Command) GetCommand(name string) *Command {
	for _, v := range cmd.children {
		if strings.EqualFold(v.name, name) {
			return v
		}
	}

	return nil
}

func (cmd *Command) NewArg(name string, desc string, param bool) *Arg {
	arg := NewArg(name, desc, param)
	arg.BindCommand(cmd)
	return arg
}

func (cmd *Command) BindArg(argv ...*Arg) {
	for _, arg := range argv {
		arg.BindCommand(cmd)
	}
}

func (cmd *Command) UnbindArg(argv ...*Arg) {
	for _, arg := range argv {
		arg.UnbindCommand(cmd)
	}
}

func (cmd *Command) GetArg(name string, param bool) *Arg {
	for _, arg := range cmd.args {
		if strings.EqualFold(arg.name, name) && arg.param == param {
			return arg
		}
	}

	// not found, may be a parameter to a parent menu
	if cmd.parent != nil {
		return cmd.parent.GetArg(name, param)
	}

	return nil
}

func (cmd *Command) hasRequiredArgs(data *Data) error {
	for _, arg := range cmd.args {
		if _, ok := data.Options[arg.name]; arg.required && !ok {
			return fmt.Errorf("Required option missing: %s", arg.name)

		}
	}

	if cmd.parent != nil {
		return cmd.parent.hasRequiredArgs(data)
	}

	return nil
}

func (cmd *Command) BindCallbackPre(handler Handler) {
	cmd.callbackspre = append(cmd.callbackspre, handler)
}

func (cmd *Command) BindCallback(handler Handler) {
	cmd.callbacks = append(cmd.callbacks, handler)
}

func (cmd *Command) runCallbacksPre(data *Data) error {
	for _, handler := range cmd.callbackspre {
		if error := handler(data); error != nil {
			return error
		}
	}

	if cmd.parent != nil {
		return cmd.parent.runCallbacksPre(data)
	}

	return nil
}

func (cmd *Command) runCallbacks(data *Data) error {
	for _, handler := range cmd.callbacks {
		if error := handler(data); error != nil {
			return error
		}
	}

	if cmd.parent != nil {
		return cmd.parent.runCallbacks(data)
	}

	return nil
}

func (cmd *Command) GetNameChain() string {
	name := cmd.name
	if cmd.parent != nil {
		parentname := cmd.parent.GetNameChain()
		if parentname != "" {
			name = parentname + " " + name
		}
	}
	return name
}

func (cmd *Command) GetNameTop() string {
	if cmd.parent != nil {
		return cmd.parent.GetNameTop()
	}

	return cmd.name
}
