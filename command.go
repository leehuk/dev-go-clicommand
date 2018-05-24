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
	// Name Name of subcommand
	Name string
	// Desc Description of subcommand
	Desc string
	// Handler Handler function subcommand calls, nil for subcommands with children
	Handler Handler
	// Parent Command object thats the parent of this one
	Parent *Command
	// Children Command objects that are children of this one
	Children []*Command
	// Options Option arguments
	Options []*Option
	// Callbackspre Callbacks to run pre-verification
	Callbackspre []Handler
	// Callbacks Callbacks to run as part of verification
	Callbacks []Handler
}

// NewCommand creates a new command, unbound to parents.  This is generally only used
// for creating the root object as after that, the func (c *Command) NewCommand()
// variant is easier, as it automatically binds the child Command.
//
// handler must be nil if this command will have its own children.
func NewCommand(name string, desc string, handler Handler) *Command {
	cmd := &Command{
		Name:    name,
		Desc:    desc,
		Handler: handler,
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
	if c.Handler != nil {
		panic(fmt.Sprintf("BindCommand() Parent has handler function set: %s", c.GetNameChain()))
	}

	c.Children = append(c.Children, cmdv...)
	for _, cmd := range cmdv {
		cmd.Parent = c
	}
}

// GetCommand finds a child Command with the given name, or nil if not found.
// name matches are case-insensitive.
func (c *Command) GetCommand(name string) *Command {
	for _, cmd := range c.Children {
		if strings.EqualFold(cmd.Name, name) {
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
	for _, option := range c.Options {
		if strings.EqualFold(option.name, name) && option.param == param {
			return option
		}
	}

	// not found, may be a parameter to a parent menu
	if c.Parent != nil {
		return c.Parent.GetOption(name, param)
	}

	return nil
}

// hasRequiredOptions iterates over all attached Option entries in the tree validating
// any marked as being required, are appropriately set.  It starts at the leaf and
// moves up towards the root.
func (c *Command) hasRequiredOptions(data *Data) error {
	for _, option := range c.Options {
		if _, ok := data.Options[option.name]; option.required && !ok {
			return fmt.Errorf("Required option missing: %s", option.name)

		}
	}

	if c.Parent != nil {
		return c.Parent.hasRequiredOptions(data)
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
	c.Callbackspre = append(c.Callbackspre, handler)
}

// BindCallback binds a validation callback, that can be used to add extra
// validation around commands and options.
//
// Callbacks are processed starting at the leaf, moving up to the root.  Only
// callbacks directly along that path are executed.
func (c *Command) BindCallback(handler Handler) {
	c.Callbacks = append(c.Callbacks, handler)
}

// runCallbacksPre runs all pre-validation callbacks, starting at the leaf
// and moving up to the root.
func (c *Command) runCallbacksPre(data *Data) error {
	for _, handler := range c.Callbackspre {
		if error := handler(data); error != nil {
			return error
		}
	}

	if c.Parent != nil {
		return c.Parent.runCallbacksPre(data)
	}

	return nil
}

// runCallbacks runs all validation callbacks, starting at the leaf and moving
// up to the root.
func (c *Command) runCallbacks(data *Data) error {
	for _, handler := range c.Callbacks {
		if error := handler(data); error != nil {
			return error
		}
	}

	if c.Parent != nil {
		return c.Parent.runCallbacks(data)
	}

	return nil
}

// GetNameChain builds a space separated string of all Command names from itself
// up to the root.
func (c *Command) GetNameChain() string {
	name := c.Name
	if c.Parent != nil {
		parentname := c.Parent.GetNameChain()
		if parentname != "" {
			name = parentname + " " + name
		}
	}
	return name
}

// GetNameTop finds the name of the root Command.
func (c *Command) GetNameTop() string {
	if c.Parent != nil {
		return c.Parent.GetNameTop()
	}

	return c.Name
}
