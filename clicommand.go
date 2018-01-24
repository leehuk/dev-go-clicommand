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
    }

    return cmd
}

func (cmd *CLICommand) AddMenu(name string, desc string, f CLICommandFunc) *CLICommand {
    subcmd := New(name, desc)
    subcmd.parent = cmd

    if f != nil {
        subcmd.f = f
    }

    cmd.children = append(cmd.children, subcmd)

    return subcmd
}

func (cmd *CLICommand) GetMenu(name string) *CLICommand {
    for _, v := range cmd.children {
        if strings.EqualFold(v.name, name) {
            return v
        }
    }

    return nil
}

func (cmd *CLICommand) GetMenuNameChain() string {
    name := cmd.name
    if cmd.parent != nil {
        parentname := cmd.parent.GetMenuNameChain()
        if parentname != "" {
            name = parentname + " " + name
        }
    }
    return name
}

func (cmd *CLICommand) GetMenuNameTop() string {
    if cmd.parent != nil {
        return cmd.parent.GetMenuNameTop()
    } else {
        return cmd.name
    }
}

func (cmd *CLICommand) AddArg(name string, desc string, param bool) {
    arg := &CLICommandArg{
        name,
        desc,
        param,
    }
    cmd.args = append(cmd.args, arg)
}

func (cmd *CLICommand) GetArg(name string, param bool) *CLICommandArg {
    for _, v := range cmd.args {
        if strings.EqualFold(v.name, name) && v.param == param {
            return v
        }
    }

    // not found, may be a parameter to a parent menu
    if cmd.parent != nil {
        return cmd.parent.GetArg(name, param)
    } else {
        return nil
    }
}

func (cmd *CLICommand) Parse() error {
    var command_ptr = cmd
    var command_data = &CLICommandData{
        options: make(map[string]string),
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
                command_data.options[argname] = argval
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
                command_data.params = os.Args[i:]
            }

            command_data.cmd = command_ptr

            command_ptr.Help(command_data)
            return nil
        // some other parameter
        } else {
            command_data.params = os.Args[i:]
            break
        }
    }

    command_data.cmd = command_ptr

    if command_ptr.f == nil {
        command_ptr.Help(command_data)
        return errors.New(fmt.Sprintf("No command specified"))
    }

    return command_ptr.f(command_data)
}

