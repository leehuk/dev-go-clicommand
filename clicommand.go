package clicommand

import(
    "errors"
    "fmt"
    "os"
)

func New(name string, desc string, parent *CLICommandMenu, f CLICommandFunc) *CLICommandMenu {
    cmd := CLICommandMenu{
        name,
        desc,
        f,
        parent,
        nil,
        nil,
    }

    if parent != nil {
        parent.children = append(parent.children, &cmd)
    }

    return &cmd
}

func NewArg(cmd *CLICommandMenu, name string, param bool, desc string) {
    arg := CLICommandArg{
        name,
        desc,
        param,
    }
    cmd.args = append(cmd.args, &arg)
}

func Parse(cmd *CLICommandMenu) error {
    var command_chain []string
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
            command_chain = append(command_chain, param)
        }
    }

    if len(command_chain) == 0 || command_chain[0] == "help" {
        Help(cmd)
    }
    fmt.Printf("FINAL: %v\n%v\n", command_chain, command_params)

    return nil
}

func GetParentName(cmd *CLICommandMenu) string {
    name := cmd.name
    if cmd.parent != nil {
        parentname := GetParentName(cmd.parent)
        if parentname != "" {
            name = parentname + " " + name
        }
    }
    return name
}


