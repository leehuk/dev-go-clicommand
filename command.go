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
// for creating the root object as after that, the func (c *Command) NewCommand()
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

// NewCommand creates a new Command and automatically binds it as a child.
//
// If the new Command will also have its own children, handler must be set to nil.  If
// the parent already has a handler set, this will panic.
func (c *Command) NewCommand(name string, desc string, handler Handler) *Command {
	cmd := NewCommand(name, desc, handler)
	c.BindCommand(cmd)
	return cmd
}

// BindCommand binds a series of subcommands as children.  Links are placed in both 
// directions, from parent -> child and child -> parent.
//
// If the parent already has a handler set, this will panic.
func (c *Command) BindCommand(cmdv ...*Command) {
	if c.handler != nil {
		panic(fmt.Sprintf("BindCommand() Parent has handler function set: %s", c.GetNameChain()))
	}

	c.children = append(c.children, cmdv...)
	for _, cmd := range cmdv {
		cmd.parent = c
	}
}

// GetCommand finds a child Command with the given name, or nil if not found.
// name matches are case-insensitive.
func (c *Command) GetCommand(name string) *Command {
	for _, cmd := range c.children {
		if strings.EqualFold(cmd.name, name) {
			return cmd
		}
	}

	return nil
}

// NewOption creates a new Option and automatically binds it as a child.
func (c *Command) NewOption(name string, desc string, param bool) *Option {
	option := NewOption(name, desc, param)
	option.BindCommand(c)
	return option
}

// BindOption binds an Option as a child.
func (c *Command) BindOption(optionv ...*Option) {
	for _, option := range optionv {
		option.BindCommand(c)
	}
}

// UnbindOption unbinds an Option so it is no longer a child.
func (c *Command) UnbindOption(optionv ...*Option) {
	for _, option := range optionv {
		option.UnbindCommand(c)
	}
}

// GetOption finds an child Option with the given name and the same parameter,
// searching the entire way up the tree to the root if necessary.
func (c *Command) GetOption(name string, param bool) *Option {
	for _, option := range c.options {
		if strings.EqualFold(option.name, name) && option.param == param {
			return option
		}
	}

	// not found, may be a parameter to a parent menu
	if c.parent != nil {
		return c.parent.GetOption(name, param)
	}

	return nil
}

// hasRequiredOptions iterates over all attached Option entries in the tree validating 
// any marked as being required, are appropriately set.  It starts at the leaf and
// moves up towards the root.
func (c *Command) hasRequiredOptions(data *Data) error {
	for _, option := range c.options {
		if _, ok := data.Options[option.name]; option.required && !ok {
			return fmt.Errorf("Required option missing: %s", option.name)

		}
	}

	if c.parent != nil {
		return c.parent.hasRequiredOptions(data)
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
func (c *Command) BindCallbackPre(handler Handler) {
	c.callbackspre = append(c.callbackspre, handler)
}

// BindCallback binds a validation callback, that can be used to add extra
// validation around commands and options.
//
// Callbacks are processed starting at the leaf, moving up to the root.  Only
// callbacks directly along that path are executed.
func (c *Command) BindCallback(handler Handler) {
	c.callbacks = append(c.callbacks, handler)
}

// runCallbacksPre runs all pre-validation callbacks, starting at the leaf
// and moving up to the root.
func (c *Command) runCallbacksPre(data *Data) error {
	for _, handler := range c.callbackspre {
		if error := handler(data); error != nil {
			return error
		}
	}

	if c.parent != nil {
		return c.parent.runCallbacksPre(data)
	}

	return nil
}

// runCallbacks runs all validation callbacks, starting at the leaf and moving
// up to the root.
func (c *Command) runCallbacks(data *Data) error {
	for _, handler := range c.callbacks {
		if error := handler(data); error != nil {
			return error
		}
	}

	if c.parent != nil {
		return c.parent.runCallbacks(data)
	}

	return nil
}

// GetNameChain() builds a space separated string of all Command names from
// itself up to the root.
func (c *Command) GetNameChain() string {
	name := c.name
	if c.parent != nil {
		parentname := c.parent.GetNameChain()
		if parentname != "" {
			name = parentname + " " + name
		}
	}
	return name
}

// GetNameTop() finds the name of the root Command.
func (c *Command) GetNameTop() string {
	if c.parent != nil {
		return c.parent.GetNameTop()
	}

	return c.name
}
