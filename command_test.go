// Copyright (C) 2018 Lee H <lee@leeh.uk>
// Licensed under the BSD 2-Clause License as found in LICENSE.txt

package clicommand

import (
	"testing"
)

func testHandlerFunc(data *Data) error {
	return nil
}

func TestNewCommand(t *testing.T) {
	var cmdRootName = "root"
	var cmdRootDesc = "root description"

	var cmdChildName = "child"

	// create our root object, and validate all struct elements are as expected
	cmdRoot := NewCommand(cmdRootName, cmdRootDesc, nil)

	if cmdRoot.name != cmdRootName {
		t.Errorf("NewCommand(cmdRoot.name); expecting %s; got %s", cmdRootName, cmdRoot.name)
	}

	if cmdRoot.desc != cmdRootDesc {
		t.Errorf("NewCommand(cmdRoot.desc); expecting %s; got %s", cmdRootDesc, cmdRoot.desc)
	}

	if cmdRoot.handler != nil {
		t.Errorf("NewCommand(cmdRoot.handler); expecting %v; got %v", nil, cmdRoot.handler)
	}

	if cmdRoot.parent != nil {
		t.Errorf("NewCommand(cmdRoot.parent); expecting %v; got %v", nil, cmdRoot.parent)
	}

	if len(cmdRoot.children) != 0 {
		t.Errorf("NewCommand(len(cmdRoot.children)); expecting 0; got %d", len(cmdRoot.children))
	}

	if len(cmdRoot.options) != 0 {
		t.Errorf("NewCommand(len(cmdRoot.options)); expecting 0; got %d", len(cmdRoot.children))
	}

	if len(cmdRoot.callbackspre) != 0 {
		t.Errorf("NewCommand(len(cmdRoot.options)); expecting 0; got %d", len(cmdRoot.children))
	}

	if len(cmdRoot.callbacks) != 0 {
		t.Errorf("NewCommand(len(cmdRoot.options)); expecting 0; got %d", len(cmdRoot.children))
	}

	// create a nested object and verify it
	cmdNest := cmdRoot.NewCommand(cmdChildName, "test", nil)

	if cmdNest.parent == nil {
		t.Errorf("NewCommand(cmdNest.parent); expecting ptr; got %v", nil)
	} else if cmdNest.parent.name != cmdRootName {
		t.Errorf("NewCommand(cmdNest.parent.name); expecting %s; got %s", cmdRootName, cmdNest.parent.name)
	}

	if len(cmdRoot.children) != 1 {
		t.Errorf("NewCommand(len(cmdRoot.children)); expecting 1; got %d", len(cmdRoot.children))
	} else {
		if cmdRoot.children[0].name != cmdChildName {
			t.Errorf("NewCommand(cmdRoot.children[0].name); expecting %s; got %s", cmdChildName, cmdRoot.children[0].name)
		}
	}
}

func TestNewCommandHandler(t *testing.T) {
	cmd := NewCommand("test", "test", testHandlerFunc)

	if cmd.handler == nil {
		t.Errorf("NewCommandHandler(handler); expecting !%v; got %v", nil, cmd.handler)
	}
}

func TestBindCommand(t *testing.T) {
	var cmdRootName = "root"
	var cmdChildName = "child"

	cmdRoot := NewCommand(cmdRootName, "root description", nil)
	cmdChild := NewCommand(cmdChildName, "child description", nil)
	cmdRoot.BindCommand(cmdChild)

	if cmdChild.parent == nil {
		t.Errorf("BindCommand(cmdChild.parent); expecting ptr; got %v", nil)
	} else if cmdChild.parent.name != cmdRootName {
		t.Errorf("BindCommand(cmdChild.parent.name); expecting %s; got %s", cmdRootName, cmdChild.parent.name)
	}

	if len(cmdRoot.children) != 1 {
		t.Errorf("BindCommand(len(cmdRoot.children)); expecting 1; got %d", len(cmdRoot.children))
	} else {
		if cmdRoot.children[0].name != cmdChildName {
			t.Errorf("BindCommand(cmdRoot.children[0].name); expecting %s; got %s", cmdChildName, cmdRoot.children[0].name)
		}
	}
}

func TestGetCommand(t *testing.T) {
	var cmdRootName = "root"
	var cmdChildName = "child"

	cmdRoot := NewCommand(cmdRootName, "root description", nil)
	cmdRoot.NewCommand(cmdChildName, "child description", nil)

	if cmdChild := cmdRoot.GetCommand(cmdChildName); cmdChild == nil {
		t.Errorf("GetCommand(cmdChildName); expecting ptr; got %v", nil)
	} else if cmdChild.name != cmdChildName {
		t.Errorf("GetCommand(cmdChildName).name; expecting %s; got %s", cmdChildName, cmdChild.name)
	}
}
