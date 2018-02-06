package clicommand

import (
	"fmt"
)

// ErrCallback Error type for when a callback has failed.
type ErrCallback struct {
	data string
}

// ErrCallbackPre Error type for when a pre-validation callback has failed.
type ErrCallbackPre struct {
	data string
}

// ErrCommandInvalid Error type for when the command line uses a subcommand
// that does not exist.
type ErrCommandInvalid struct {
	data string
}

// ErrCommandMissing Error type for when the command line has ended with a
// parent command with no handler, meaning one of its children needed to be
// chosen instead.
type ErrCommandMissing struct{}

// ErrOptionMissing Error type for when a required option is missing.
type ErrOptionMissing struct {
	data string
}

// ErrOptionMissingParam Error type for when the command line contains an
// option that requires a parameter, but one is not specified
type ErrOptionMissingParam struct {
	data string
}

// ErrOptionUnknown Error type for when the command line contains an option,
// that is not defined in the command tree.
type ErrOptionUnknown struct {
	data string
}

func (e *ErrCallback) Error() string {
	return fmt.Sprintf("Callback error: %s", e.data)
}

func (e *ErrCallbackPre) Error() string {
	return fmt.Sprintf("CallbackPre error: %s", e.data)
}

func (e *ErrCommandInvalid) Error() string {
	return fmt.Sprintf("Invalid subcommand: %s", e.data)
}

func (e *ErrCommandMissing) Error() string {
	return fmt.Sprintf("No subcommand specified")
}

func (e *ErrOptionMissing) Error() string {
	return fmt.Sprintf("Required option missing: %s", e.data)
}

func (e *ErrOptionMissingParam) Error() string {
	return fmt.Sprintf("Missing parameter to option: %s", e.data)
}

func (e *ErrOptionUnknown) Error() string {
	return fmt.Sprintf("Unknown option: %s", e.data)
}
