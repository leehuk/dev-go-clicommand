package clicommand

import(
    "fmt"
)

func (cmd *CLICommand) Help() {
    fmt.Printf("\n")
    fmt.Printf("  %s\n", cmd.GetMenuNameChain())
    fmt.Printf("  %s\n", cmd.desc)
    fmt.Printf("\n")

    if len(cmd.children) > 0 {
        fmt.Printf("  Available subcommands:\n")
        for _, v := range cmd.children {
            fmt.Printf("    %-12s %s\n", v.name, v.desc)
        }
        fmt.Printf("\n")

        cmd.HelpOptions()

        fmt.Printf("  For a command overview run 'gh help', 'gh <command> help',\n")
        fmt.Printf("  'gh <command> <subcommand> help', etc\n")
        fmt.Printf("\n")
    }
}

func (cmd *CLICommand) HelpOptions() {
    fmt.Printf("  %s options:\n", cmd.GetMenuNameChain())
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
        cmd.parent.HelpOptions()
    }
}
