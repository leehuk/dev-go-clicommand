// Copyright (C) 2018 Lee H <lee@leeh.uk>
// Licensed under the BSD 2-Clause License as found in LICENSE.txt

package clicommand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cmdRootName  string = "root"
	cmdRootDesc  string = "root description"
	cmdChildName string = "child"
	cmdChildDesc string = "child description"

	optionName  string = "option"
	optionDesc  string = "option description"
	optionParam        = false
)

func testHandlerFunc(data *Data) error {
	return nil
}

func newCommandRoot(handler Handler) *Command {
	return NewCommand(cmdRootName, cmdRootDesc, handler)
}

func newCommandChild(handler Handler) *Command {
	return NewCommand(cmdChildName, cmdChildDesc, handler)
}

func (cmd *Command) newCommandChild(handler Handler) *Command {
	return cmd.NewCommand(cmdChildName, cmdChildDesc, handler)
}

func newOption() *Option {
	return NewOption(optionName, optionDesc, optionParam)
}

func (cmd *Command) newOption() *Option {
	return cmd.NewOption(optionName, optionDesc, optionParam)
}

// NewCommand testing, both unbound and bound forms, with/without handlers
func TestNewCommand(t *testing.T) {
	assert := assert.New(t)

	// create our root object, and validate all struct elements are as expected
	cmdRoot := newCommandRoot(nil)
	assert.Equal(cmdRootName, cmdRoot.Name, "command.Name error")
	assert.Equal(cmdRootDesc, cmdRoot.Desc, "command.Desc error")
	assert.Nil(cmdRoot.Handler)
	assert.Nil(cmdRoot.Parent)
	assert.Empty(cmdRoot.Children)
	assert.Empty(cmdRoot.Options)
	assert.Empty(cmdRoot.Callbackspre)
	assert.Empty(cmdRoot.Callbacks)

	// create a nested object and verify it
	cmdChild := cmdRoot.newCommandChild(nil)

	if assert.NotNil(cmdChild.Parent) {
		assert.Equal(cmdRootName, cmdChild.Parent.Name, "Command.parent.name error")
	}

	if assert.Len(cmdRoot.Children, 1) {
		assert.Equal(cmdChildName, cmdRoot.Children[0].Name, "Command.children[0].name")
	}

	cmdHandler := newCommandRoot(testHandlerFunc)
	assert.NotNil(cmdHandler.Handler)
}

// BindCommand testing, create parent/child objects and bind them
func TestBindCommand(t *testing.T) {
	assert := assert.New(t)

	cmdRoot := newCommandRoot(nil)
	cmdChild := newCommandChild(nil)
	cmdRoot.BindCommand(cmdChild)

	if assert.NotNil(cmdChild.Parent) {
		assert.Equal(cmdRootName, cmdChild.Parent.Name)
	}

	if assert.Len(cmdRoot.Children, 1) {
		assert.Equal(cmdChildName, cmdRoot.Children[0].Name)
	}
}

// GetCommand testing, create bound parent/child and validate we can find
// the child
func TestGetCommand(t *testing.T) {
	assert := assert.New(t)

	cmdRoot := newCommandRoot(nil)
	cmdRoot.newCommandChild(nil)

	cmdChild := cmdRoot.GetCommand(cmdChildName)

	if assert.NotNil(cmdChild) {
		assert.Equal(cmdChildName, cmdChild.Name)
	}
}

// NewOption testing, basic validation it is bound, most testing in option_test.go
func TestNewOption(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	option := cmd.newOption()

	if assert.Len(option.parents, 1) {
		assert.Len(cmd.Options, 1)
	}
}

// BindOption testing, validation it is bound
func TestBindOption(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	option := newOption()
	cmd.BindOption(option)

	assert.Len(option.parents, 1)
	assert.Len(cmd.Options, 1)
}

// UnbindOption, validate we bind and then unbind
func TestUnbindOption(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	option := cmd.newOption()

	// first verify its bound
	if assert.Len(option.parents, 1) {
		cmd.UnbindOption(option)
		assert.Empty(option.parents)
		assert.Empty(cmd.Options)
	}

}

// GetOption, create bound cmd/option and validate we can find it,
// then add child and validate we can find the parent option through it
func TestGetOption(t *testing.T) {
	assert := assert.New(t)

	cmdRoot := newCommandRoot(nil)
	cmdRoot.newOption()

	option := cmdRoot.GetOption(optionName, optionParam)
	assert.NotNil(option)

	cmdChild := cmdRoot.newCommandChild(nil)
	option = cmdChild.GetOption(optionName, optionParam)
	assert.NotNil(option)
}

// BindCallbackPre, create simple callback and validate its added
func TestBindCallbackPre(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	cmd.BindCallbackPre(testHandlerFunc)

	assert.Len(cmd.Callbackspre, 1)
}

// BindCallback, create simple callback and validate its added
func TestBindCallback(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	cmd.BindCallback(testHandlerFunc)

	assert.Len(cmd.Callbacks, 1)
}

// TestGetNameChain, create parent/child commands, and validate name
func TestGetNameChain(t *testing.T) {
	assert := assert.New(t)

	cmdRoot := newCommandRoot(nil)
	cmdChild := cmdRoot.newCommandChild(nil)

	assert.Equal(cmdRootName+" "+cmdChildName, cmdChild.GetNameChain(), "Command.GetNameChild()")
}

func TestGetNameParent(t *testing.T) {
	assert := assert.New(t)

	cmdRoot := newCommandRoot(nil)
	cmdChild := cmdRoot.newCommandChild(nil)

	assert.Equal(cmdRootName, cmdChild.GetNameTop(), "Command.GetNameTop()")
}
