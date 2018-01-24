package clicommand

import(
    "errors"
    "fmt"
    "os"
    "strings"
)

func New(name string, desc string) *CLICommand {
    cmd := &CLICommand{
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

func (cmd *CLICommand) Parse() error {
    var command_ptr = cmd
    var command_data = &CLICommandData{
        Options: make(map[string]string),
    }

    // no parameters given, display overall help
    if len(os.Args) <= 1 {
        command_ptr.Help(nil)
        return nil
    }

    for i := 1; i < len(os.Args); i++ {
        arg := os.Args[i]

        // option argument
        if len(arg) >= 1 && arg[:1] == "-" {
            var argname string
            var argval string
            var argparam bool

            // ensure we do not have an option with no name
            if len(arg) == 1 && arg[:1] == "-" || len(arg) == 2 && arg[:2] == "--" {
                return errors.New(fmt.Sprintf("Invalid option: %s", arg))
            }

            // option with parameter: "--xyz"
            if arg[:2] == "--" {
                // ensure we have a parameter
                if i+1 >= len(os.Args) {
                    return errors.New(fmt.Sprintf("Missing parameter to option: %s", arg))
                }

                argname = arg[2:]
                argval = os.Args[i+1]
                argparam = true

                // next arg was an option to this param, skip its parsing
                i++
                // option without parameter: "-xyz"
            } else {
                argname = arg[1:]
                argval = ""
                argparam = false
            }

            if subarg := command_ptr.GetArg(argname, argparam); subarg != nil {
                command_data.Options[argname] = argval
            } else {
                return errors.New(fmt.Sprintf("Unknown option: %s", arg))
            }
        // sub-menu
        } else if subcmd := command_ptr.GetMenu(arg); subcmd != nil {
            // repoint our pointer to this sub-menu and continue parsing
            command_ptr = subcmd
        // help command as sub-menu.  This calls directly out to Help() on the current
        // sub-command object, then returns.
        } else if strings.EqualFold(arg, "help") {
            // take any remaining fields as parameters
            if len(os.Args) >= i {
                i++
                command_data.Params = os.Args[i:]
            }

            command_data.Cmd = command_ptr
            command_ptr.Help(command_data)

            return nil
        // some other parameter
        } else {
            command_data.Params = os.Args[i:]
            break
        }
    }

    command_data.Cmd = command_ptr

    if command_ptr.f == nil {
        command_ptr.Help(command_data)
        return errors.New(fmt.Sprintf("No command specified"))
    }

    if e := command_ptr.RunCallbacks(command_data); e != nil {
        return e
    }

    return command_ptr.f(command_data)
}

