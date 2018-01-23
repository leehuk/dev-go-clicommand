package clicommand

import(
    "errors"
    "fmt"
    "os"
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

func (cmd *CLICommand) AddArg(name string, desc string, param bool) {
    arg := &CLICommandArg{
        name,
        desc,
        param,
    }
    cmd.args = append(cmd.args, arg)
}

func (cmd *CLICommand) Parse() error {
    var command_params = make(map[string]string)

    for i := 1; i < len(os.Args); i++ {
        param := os.Args[i]

        // double dash options have parameters
        if len(param) >= 3 && param[:2] == "--" {
            if i + 1 >= len(os.Args) {
                return errors.New(fmt.Sprintf("Missing parameter value for option: %s", param))
            }

            pkey := param[2:]
            pvalue := os.Args[i+1]

            command_params[pkey] = pvalue

            // skip next parameter
            i++
        // single dash parameter, no value
        } else if len(param) >= 2 && param[:1] == "-" {
            pkey := param[1:]
            command_params[pkey] = ""
        } else {
        }

    }

    fmt.Printf("FINAL: %v\n%v\n", command_params)

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


