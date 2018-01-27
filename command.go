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
	// options Option arguments
	options      []*Option
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

// NewOption creates a new Option and automatically binds it as a child of self.
func (self *Command) NewOption(name string, desc string, param bool) *Option {
	option := NewOption(name, desc, param)
	option.BindCommand(self)
	return option
}

// BindOption binds an Option as a child of self.
func (self *Command) BindOption(optionv ...*Option) {
	for _, option := range optionv {
		option.BindCommand(self)
	}
}

// UnbindOption unbinds an Option so it is no longer a child of self.
func (self *Command) UnbindOption(optionv ...*Option) {
	for _, option := range optionv {
		option.UnbindCommand(self)
	}
}

// GetOption finds a child Option of self with the given name and the same parameter
// type, or nil if not found.
func (self *Command) GetOption(name string, param bool) *Option {
	for _, option := range self.options {
		if strings.EqualFold(option.name, name) && option.param == param {
			return option
		}
	}

	// not found, may be a parameter to a parent menu
	if self.parent != nil {
		return self.parent.GetOption(name, param)
	}

	return nil
}

// hasRequiredOptions iterates over all attached Option entries in the tree validating 
// any marked as being required, are appropriately set.  It starts at the leaf and
// moves up towards the root.
func (self *Command) hasRequiredOptions(data *Data) error {
	for _, option := range self.options {
		if _, ok := data.Options[option.name]; option.required && !ok {
			return fmt.Errorf("Required option missing: %s", option.name)

		}
	}

	if self.parent != nil {
		return self.parent.hasRequiredOptions(data)
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
