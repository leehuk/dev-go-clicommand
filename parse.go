package clicommand

import (
	"fmt"
	"os"
	"strings"
)

// Parse parses the command line from os.Args under the supplied command tree, then
// acts accordingly based on the results.
//
// Everything specified on the command line either guides us towards the correct
// command to run, is used as an option, or becomes a generic parameter.
//
// Once parsing is complete, pre callbacks are made, then we either proceed to
// display internal help information if requested, or we perform internal
// verification, then call the validation callbacks, then finally if everything
// is ok call the wanted Handler.
//
// The parsing will steal the arg "help" if it detects it as the first unknown
// parameter, allowing for easy access to the available commands and options.
func (c *Command) Parse() error {
	var commandPtr = c
	var commandData = &Data{
		Cmd:     c,
		Options: make(map[string]string),
	}

	var param_parsing = false
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if len(arg) >= 1 && arg[:1] == "-" {
			// option argument
			var optionname string
			var optionval string
			var optionparam bool

			// ensure we do not have an option with no name
			if len(arg) == 1 && arg[:1] == "-" || len(arg) == 2 && arg[:2] == "--" {
				return fmt.Errorf("Invalid option: %s", arg)
			}

			if arg[:2] == "--" {
				// option with parameter: "--xyz"

				// ensure we have a parameter
				if i+1 >= len(os.Args) {
					return fmt.Errorf("Missing parameter to option: %s", arg)
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
				return fmt.Errorf("Unknown option: %s", arg)
			}
		} else if param_parsing {
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
		} else {
			// some other parameter
			commandData.Params = append(commandData.Params, os.Args[i])
			// At this point, we only parse things into options or parameters
			param_parsing = true
		}
	}

	// no subcommand specified
	if commandPtr.handler == nil {
		// dont error if we're at the root level
		if commandPtr == c {
			helpUsage(commandData)
			return nil
		}

		return helpError(commandData, fmt.Errorf("No subcommand specified"))
	}

	if e := commandPtr.runCallbacksPre(commandData); e != nil {
		return helpError(commandData, e)
	}

	if commandPtr != cmdHelp {
		if e := commandPtr.hasRequiredOptions(commandData); e != nil {
			return helpError(commandData, e)
		}

		if e := commandPtr.runCallbacks(commandData); e != nil {
			return helpError(commandData, e)
		}
	}

	return commandPtr.handler(commandData)
}
