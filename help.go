package clicommand

import(
    "fmt"
)

func Help(cmd *Command) {
    fmt.Printf("\n")
    fmt.Printf("  %s\n", GetParentName(cmd))
    fmt.Printf("  %s\n", cmd.desc)
    fmt.Printf("\n")

    if len(cmd.children) > 0 {
        fmt.Printf("  Available subcommands:\n")
        for _, v := range cmd.children {
            fmt.Printf("    %-12s %s\n", v.name, v.desc)
        }
        fmt.Printf("\n")

        HelpOptions(cmd)

        fmt.Printf("  For a command overview run 'gh help', 'gh <command> help',\n")
        fmt.Printf("  'gh <command> <subcommand> help', etc\n")
        fmt.Printf("\n")
    }
}

func HelpOptions(cmd *Command) {
    fmt.Printf("  %s options:\n", GetParentName(cmd))
    for _, arg := range cmd.args {
        var prefix string

        if arg.param {
            prefix = "--"
        } else {
            prefix = "-"
        }

        fmt.Printf("    %2s%-14s %s\n", prefix, arg.name, arg.desc)
    }

    fmt.Printf("\n")

    if cmd.parent != nil {
        HelpOptions(cmd.parent)
    }
}
