// Copyright (C) 2018 Lee H <lee@leeh.uk>
// Licensed under the BSD 2-Clause License as found in LICENSE.txt

package clicommand

import (
	"fmt"
	"strings"
)

// A Command represents a command of the cli program.  These are chained into a tree 
// with links in both child and parent directions
//
// Each child can itself be a parent, but only leaf nodes (children with no children)
// can have handler functions, e.g.
//  clicommand -->
//    clicommand api -->
//      clicommand api get ==> handler
//      clicommand api delete ==> handler
type Command struct {
	// name Name of subcommand
	name         string
	// desc Description of subcommand
	desc         string
	// handler Handler function subcommand calls, nil for subcommands with children
	handler      Handler
	// parent Command object thats the parent of this one
	parent       *Command
	// children Command objects that are children of this one
	children     []*Command
	// args Option arguments
	args         []*Arg
	// callbackspre Callbacks to run pre-verification
	callbackspre []Handler
	// callbacks Callbacks to run as part of verification
	callbacks    []Handler
}

// NewCommand creates a new command, unbound to parents.  This is generally only used
// for creating the root object as after that, the func (self *Command) NewCommand()
// variant is easier, as it automatically binds the child Command.
//
// handler must be nil if this command will have its own children.
func NewCommand(name string, desc string, handler Handler) *Command {
	cmd := &Command{
		name:    name,
		desc:    desc,
		handler: handler,
	}

	return cmd
}

// NewCommand creates a new Command and automatically binds it as a child of self.
//
// If the new Command will also have its own children, handler must be set to nil.  If
// self already has a handler set, this will panic.
func (self *Command) NewCommand(name string, desc string, handler Handler) *Command {
	cmd := NewCommand(name, desc, handler)
	self.BindCommand(cmd)
	return cmd
}

// BindCommand binds a series of subcommands as children of self.  Links are placed
// in both directions, from parent -> child and child -> parent.
//
// If self already has a handler set, this will panic.
func (self *Command) BindCommand(cmdv ...*Command) {
	if self.handler != nil {
		panic(fmt.Sprintf("BindCommand() Parent has handler function set: %s", self.GetNameChain()))
	}

	self.children = append(self.children, cmdv...)
	for _, cmd := range cmdv {
		cmd.parent = self
	}
}

// GetCommand finds a child Command of self with the given name, or nil if not found.
// name matches are case-insensitive.
func (self *Command) GetCommand(name string) *Command {
	for _, cmd := range self.children {
		if strings.EqualFold(cmd.name, name) {
			return cmd
		}
	}

	return nil
}

// NewArg creates a new Arg and automatically binds it as a child of self.
func (self *Command) NewArg(name string, desc string, param bool) *Arg {
	arg := NewArg(name, desc, param)
	arg.BindCommand(self)
	return arg
}

// BindArg binds an Arg as a child of self.
func (self *Command) BindArg(argv ...*Arg) {
	for _, arg := range argv {
		arg.BindCommand(self)
	}
}

// UnbindArg unbinds an Arg so it is no longer a child of self.
func (self *Command) UnbindArg(argv ...*Arg) {
	for _, arg := range argv {
		arg.UnbindCommand(self)
	}
}

// GetArg finds a child Arg of self with the given name and the same parameter
// type, or nil if not found.
func (self *Command) GetArg(name string, param bool) *Arg {
	for _, arg := range self.args {
		if strings.EqualFold(arg.name, name) && arg.param == param {
			return arg
		}
	}

	// not found, may be a parameter to a parent menu
	if self.parent != nil {
		return self.parent.GetArg(name, param)
	}

	return nil
}

// hasRequiredArgs iterates over all attached Arg entries in the tree validating 
// any marked as being required, are appropriately set.  It starts at the leaf and
// moves up towards the root.
func (self *Command) hasRequiredArgs(data *Data) error {
	for _, arg := range self.args {
		if _, ok := data.Options[arg.name]; arg.required && !ok {
			return fmt.Errorf("Required option missing: %s", arg.name)

		}
	}

	if self.parent != nil {
		return self.parent.hasRequiredArgs(data)
	}

	return nil
}

// BindCallbackPre binds a pre-validation callback, that can be used to alter
// the user-provided options and Command tree prior to validation.  This can be
// useful for things like translating environment variables into options, or
// making options required/not-required for certain commands.
//
// Callbacks are processed starting at the leaf, moving up to the root. Only
// callbacks directly on that path are executed.
func (self *Command) BindCallbackPre(handler Handler) {
	self.callbackspre = append(self.callbackspre, handler)
}

// BindCallback binds a validation callback, that can be used to add extra
// validation around commands and options.
//
// Callbacks are processed starting at the leaf, moving up to the root.  Only
// callbacks directly along that path are executed.
func (self *Command) BindCallback(handler Handler) {
	self.callbacks = append(self.callbacks, handler)
}

// runCallbacksPre runs all pre-validation callbacks, starting at the leaf
// and moving up to the root.
func (self *Command) runCallbacksPre(data *Data) error {
	for _, handler := range self.callbackspre {
		if error := handler(data); error != nil {
			return error
		}
	}

	if self.parent != nil {
		return self.parent.runCallbacksPre(data)
	}

	return nil
}

// runCallbacks runs all validation callbacks, starting at the leaf and moving
// up to the root.
func (self *Command) runCallbacks(data *Data) error {
	for _, handler := range self.callbacks {
		if error := handler(data); error != nil {
			return error
		}
	}

	if self.parent != nil {
		return self.parent.runCallbacks(data)
	}

	return nil
}

// GetNameChain() builds a space separated string of all Command names from
// self up to the root.
func (self *Command) GetNameChain() string {
	name := self.name
	if self.parent != nil {
		parentname := self.parent.GetNameChain()
		if parentname != "" {
			name = parentname + " " + name
		}
	}
	return name
}

// GetNameTop() finds the name of the root Command.
func (self *Command) GetNameTop() string {
	if self.parent != nil {
		return self.parent.GetNameTop()
	}

	return self.name
}
