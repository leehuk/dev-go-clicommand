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
    var command_params = make(map[string]string)
    var command_ptr = cmd

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

            // option with parameter "--xyz"
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
            // option without parameter "-xyz"
            } else {
                argname = arg[1:]
                argval = ""
                argparam = false
            }

            if subarg := command_ptr.GetArg(argname, argparam); subarg != nil {
                command_params[argname] = argval
            } else {
                return errors.New(fmt.Sprintf("Unknown option: %s", arg))
            }
        // sub-menu
        } else if subcmd := command_ptr.GetMenu(arg); subcmd != nil {
            command_ptr = subcmd
        } else if strings.EqualFold(arg, "help") {
            command_ptr.Help()
        } else {
        }
    }

    if command_ptr == cmd {
        command_ptr.Help()
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


