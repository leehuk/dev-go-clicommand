package clicommand

import (
	"os"
	"strings"
)

// Parse parses the command line from os.Args under the supplied command tree, then
// acts accordingly based on the results.
//
// Everything specified on the command line is either a subcommand, option or a
// generic parameter.
//
// Once parsing is complete, pre callbacks are made, then we either proceed to
// display internal help information if requested, or we perform internal
// verification, then call the validation callbacks, then finally if everything
// is ok call the wanted Handler.
//
// The parsing will steal the arg "help" if it detects it as the first unknown
// parameter, allowing for easy access to the available commands and options.
//
// If parsing is not ok, it will return one of several internal error types.
func (c *Command) Parse() error {
	var commandPtr = c
	var commandData = &Data{
		Cmd:     c,
		Options: make(map[string]string),
	}

	var paramParsing = false
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if len(arg) >= 1 && arg[:1] == "-" {
			// option argument
			var optionname string
			var optionval string
			var optionparam bool

			// ensure we do not have an option with no name
			if len(arg) == 1 && arg[:1] == "-" || len(arg) == 2 && arg[:2] == "--" {
				return &ErrOptionUnknown{arg}
			}

			if arg[:2] == "--" {
				// option with parameter: "--xyz"

				// ensure we have a parameter
				if i+1 >= len(os.Args) {
					return &ErrOptionMissingParam{arg}
				}

				optionname = arg[2:]
				optionval = os.Args[i+1]
				optionparam = true

				// next arg was an option to this param, skip its parsing
				i++
			} else {
				// option without parameter: "-xyz"

				optionname = arg[1:]
				optionval = ""
				optionparam = false
			}

			if subarg := commandPtr.GetOption(optionname, optionparam); subarg != nil {
				commandData.Options[optionname] = optionval
			} else {
				return &ErrOptionUnknown{arg}
			}
		} else if paramParsing {
			// parameter parsing
			commandData.Params = append(commandData.Params, os.Args[i])
		} else if subcmd := commandPtr.GetCommand(arg); subcmd != nil {
			// sub-menu

			// repoint our pointer to this sub-menu and continue parsing
			commandPtr = subcmd
			commandData.Cmd = commandPtr
		} else if strings.EqualFold(arg, "help") {
			// help command as sub-menu

			// take any remaining fields as parameters
			if len(os.Args) >= i {
				commandData.Params = os.Args[i+1:]
				i = len(os.Args)
			}

			// we now want to call out to help on a dummy command object, but preserving
			// Cmd as our current position down the menu structure
			commandData.Cmd = commandPtr
			cmdHelp.parent = commandPtr
			commandPtr = cmdHelp
		} else if commandPtr.handler == nil {
			// we're in a parent menu, so this cant be a parameter -- but the next argument
			// is not a valid subcommand.
			return &ErrCommandInvalid{arg}
		} else {
			// we've now reached a child menu, and all that remains are parameters and options
			commandData.Params = append(commandData.Params, os.Args[i])
			paramParsing = true
		}
	}

	// no subcommand specified
	if commandPtr.handler == nil {
		// dont error if we're at the root level
		if commandPtr == c {
			helpUsage(commandData)
			return nil
		}

		return &ErrCommandMissing{}
	}

	if e := commandPtr.runCallbacksPre(commandData); e != nil {
		return &ErrCallbackPre{e.Error()}
	}

	if commandPtr != cmdHelp {
		if e := commandPtr.hasRequiredOptions(commandData); e != nil {
			return &ErrOptionMissing{e.Error()}
		}

		if e := commandPtr.runCallbacks(commandData); e != nil {
			return &ErrCallback{e.Error()}
		}
	}

	return commandPtr.handler(commandData)
}
