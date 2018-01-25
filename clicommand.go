package clicommand

import (
	"fmt"
	"os"
	"strings"
)

func New(name string, desc string) *Command {
	cmd := &Command{
		name,
		desc,
		nil,
		nil,
		nil,
		nil,
		nil,
	}

	return cmd
}

func (cmd *Command) Parse() error {
	var commandPtr = cmd
	var commandData = &Data{
		Options: make(map[string]string),
	}

	// no parameters given, display overall help
	if len(os.Args) <= 1 {
		commandPtr.Help(nil)
		return nil
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
		} else if strings.EqualFold(arg, "help") {
			// help command as sub-menu.  This calls directly out to Help() on the current
			// sub-command object, then returns.

			// take any remaining fields as parameters
			if len(os.Args) >= i {
				i++
				commandData.Params = os.Args[i:]
			}

			commandData.Cmd = commandPtr
			commandPtr.Help(commandData)

			return nil
		} else {
			// some other parameter
			commandData.Params = os.Args[i:]
			break
		}
	}

	commandData.Cmd = commandPtr

	if commandPtr.handler == nil {
		commandPtr.Help(commandData)
		return fmt.Errorf("No command specified")
	}

	if e := commandPtr.RunCallbacks(commandData); e != nil {
		return e
	}

	return commandPtr.handler(commandData)
}
