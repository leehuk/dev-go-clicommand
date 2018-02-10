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
	assert.Equal(cmdRootName, cmdRoot.name, "command.Name error")
	assert.Equal(cmdRootDesc, cmdRoot.desc, "command.Desc error")
	assert.Nil(cmdRoot.handler)
	assert.Nil(cmdRoot.parent)
	assert.Empty(cmdRoot.children)
	assert.Empty(cmdRoot.options)
	assert.Empty(cmdRoot.callbackspre)
	assert.Empty(cmdRoot.callbacks)

	// create a nested object and verify it
	cmdChild := cmdRoot.newCommandChild(nil)

	if assert.NotNil(cmdChild.parent) {
		assert.Equal(cmdRootName, cmdChild.parent.name, "Command.parent.name error")
	}

	if assert.Len(cmdRoot.children, 1) {
		assert.Equal(cmdChildName, cmdRoot.children[0].name, "Command.children[0].name")
	}

	cmdHandler := newCommandRoot(testHandlerFunc)
	assert.NotNil(cmdHandler.handler)
}

// BindCommand testing, create parent/child objects and bind them
func TestBindCommand(t *testing.T) {
	assert := assert.New(t)

	cmdRoot := newCommandRoot(nil)
	cmdChild := newCommandChild(nil)
	cmdRoot.BindCommand(cmdChild)

	if assert.NotNil(cmdChild.parent) {
		assert.Equal(cmdRootName, cmdChild.parent.name)
	}

	if assert.Len(cmdRoot.children, 1) {
		assert.Equal(cmdChildName, cmdRoot.children[0].name)
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
		assert.Equal(cmdChildName, cmdChild.name)
	}
}

// NewOption testing, basic validation it is bound, most testing in option_test.go
func TestNewOption(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	option := cmd.newOption()

	if assert.Len(option.parents, 1) {
		assert.Len(cmd.options, 1)
	}
}

// BindOption testing, validation it is bound
func TestBindOption(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	option := newOption()
	cmd.BindOption(option)

	assert.Len(option.parents, 1)
	assert.Len(cmd.options, 1)
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
		assert.Empty(cmd.options)
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

	assert.Len(cmd.callbackspre, 1)
}

// BindCallback, create simple callback and validate its added
func TestBindCallback(t *testing.T) {
	assert := assert.New(t)

	cmd := newCommandRoot(nil)
	cmd.BindCallback(testHandlerFunc)

	assert.Len(cmd.callbacks, 1)
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
