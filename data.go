// Copyright (C) 2018 Lee H <lee@leeh.uk>
// Licensed under the BSD 2-Clause License as found in LICENSE.txt

package clicommand

// The Data structure is passed to all Handler functions called as a result
// of a given Command being run.  This structure is also passed to any
// registered callbacks during the parsing stage.
type Data struct {
	// Cmd is a pointer to the Command object which has triggered the Handler
	// to be called.
	//
	// If the Handler is a registered callback, the pointer is to the Command
	// object that is about to be executed.
	Cmd *Command

	// Options is a map of options supplied to the command.  The key is the
	// option selected by the user, with the value being the parameter supplied
	// to that option.
	//
	// For Option objects which do not take parameters, the value is an empty
	// string.
	//
	// E.g.:
	//   ./clicommand ... --foo bar ... -q ...
	// Becomes:
	//   Options["foo"] = "bar"
	//   Options["q"] = ""
	//
	// Callbacks may modify the Options directly, the parser then sees these
	// as if they were directly supplied by the user.  Callbacks of this type
	// should generally be bound via BindCallbackPre(), allowing verification
	// to be performed as normal within standard callbacks.
	Options map[string]string

	// Params is an array of additional parameters the parser did not recognise.
	// Effectively, when the parser finds a non-option argument which doesnt
	// match any more commands, the remaining non-option fields become
	// parameters.
	Params []string
}
