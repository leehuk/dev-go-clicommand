package clicommand

import (
	"fmt"
	"os"
)

var (
	cmdHelp = &Command{
		handler: helpUsage,
	}
)

func helpError(data *Data, err error) error {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "For help information, run: %s help\n", data.Cmd.GetNameChain())
	return fmt.Errorf("%s", err)
}

func helpUsage(data *Data) error {
	cmd := data.Cmd

	fmt.Printf("\n")
	fmt.Printf("%s\n", cmd.GetNameChain())
	fmt.Printf("%s\n", cmd.desc)
	fmt.Printf("\n")

	helpOptionsRecurseRev(cmd)

	if len(cmd.children) > 0 {
		fmt.Printf("Available subcommands:\n")
		for _, v := range cmd.children {
			fmt.Printf("  %-12s %s\n", v.name, v.desc)
		}
		fmt.Printf("\n")
	}

	if cmd.handler == nil {
		fmt.Printf("For help information run:\n")
		fmt.Printf("  '%s help' .. '%s <commands>* help' .. '%s [commands]* help [subcommand]*'\n",
			cmd.GetNameTop(), cmd.GetNameTop(), cmd.GetNameTop())
		fmt.Printf("\n")
	}

	return nil
}

func helpOptionsRecurseRev(cmd *Command) {
	if cmd.parent != nil {
		helpOptionsRecurseRev(cmd.parent)
	}

	helpOptions(cmd)
}

func helpOptions(cmd *Command) {
	if len(cmd.args) == 0 {
		return
	}

	fmt.Printf("%s options:\n", cmd.GetNameChain())
	for _, arg := range cmd.args {
		var opttype string
		var optsuffix string
		var descprefix string

		if arg.param {
			opttype += "--"
			optsuffix += " <arg>"
		} else {
			opttype += "-"
		}

		if arg.required {
			descprefix += "Required: "
		}

		fmt.Printf("  %2s%-20s %s\n", opttype, arg.name+optsuffix, descprefix+arg.desc)
	}

	fmt.Printf("\n")
}
