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
	helpOutput(data, true)

	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "For help information, run: %s help\n", data.Cmd.GetNameChain())

	return err
}

func helpUsage(data *Data) error {
	helpOutput(data, false)
	return nil
}

func helpOutput(data *Data, stderr bool) {
	out := os.Stdout
	if stderr {
		out = os.Stderr
	}

	cmd := data.Cmd

	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "%s\n", cmd.GetNameChain())
	fmt.Fprintf(out, "%s\n", cmd.desc)
	fmt.Fprintf(out, "\n")

	helpOptionsRecurseRev(cmd)

	if len(cmd.children) > 0 {
		fmt.Fprintf(out, "Available subcommands:\n")
		for _, v := range cmd.children {
			fmt.Fprintf(out, "  %-12s %s\n", v.name, v.desc)
		}
		fmt.Fprintf(out, "\n")
	}

	if cmd.handler == nil {
		fmt.Fprintf(out, "For help information run:\n")
		fmt.Fprintf(out, "  '%s help' .. '%s <commands>* help' .. '%s [commands]* help [subcommand]*'\n",
			cmd.GetNameTop(), cmd.GetNameTop(), cmd.GetNameTop())
		fmt.Fprintf(out, "\n")
	}
}

func helpOptionsRecurseRev(cmd *Command) {
	if cmd.parent != nil {
		helpOptionsRecurseRev(cmd.parent)
	}

	helpOptions(cmd)
}

func helpOptions(cmd *Command) {
	if len(cmd.options) == 0 {
		return
	}

	fmt.Printf("%s options:\n", cmd.GetNameChain())
	for _, option := range cmd.options {
		var opttype string
		var optsuffix string
		var descprefix string

		if option.param {
			opttype += "--"
			optsuffix += " <arg>"
		} else {
			opttype += "-"
		}

		if option.required {
			descprefix += "Required: "
		}

		fmt.Printf("  %2s%-20s %s\n", opttype, option.name+optsuffix, descprefix+option.desc)
	}

	fmt.Printf("\n")
}
