package clicommand

import (
	"fmt"
	"os"
	"strings"
)

func (cmd *Command) Parse() error {
	var commandPtr = cmd
	var commandData = &Data{
		Cmd:     cmd,
		Options: make(map[string]string),
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if len(arg) >= 1 && arg[:1] == "-" {
			// option argument
			var argname string
			var argval string
			var argparam bool

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

				argname = arg[2:]
				argval = os.Args[i+1]
				argparam = true

				// next arg was an option to this param, skip its parsing
				i++
			} else {
				// option without parameter: "-xyz"

				argname = arg[1:]
				argval = ""
				argparam = false
			}

			if subarg := commandPtr.GetArg(argname, argparam); subarg != nil {
				commandData.Options[argname] = argval
			} else {
				return fmt.Errorf("Unknown option: %s", arg)
			}
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
		}
	}

	// no subcommand specified
	if commandPtr.handler == nil {
		// dont error if we're at the root level
		if commandPtr == cmd {
			helpUsage(commandData)
			return nil
		} else {
			return helpError(commandData, fmt.Errorf("No subcommand specified"))
		}
	}

	if e := commandPtr.runCallbacksPre(commandData); e != nil {
		return helpError(commandData, e)
	}

	if commandPtr != cmdHelp {
		if e := commandPtr.hasRequiredArgs(commandData); e != nil {
			return helpError(commandData, e)
		}

		if e := commandPtr.runCallbacks(commandData); e != nil {
			return helpError(commandData, e)
		}
	}

	return commandPtr.handler(commandData)
}
