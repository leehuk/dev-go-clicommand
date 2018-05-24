package clicommand

import (
	"fmt"
	"os"
	"strings"
)

var (
	cmdHelp = &Command{
		Handler: helpUsage,
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
	fmt.Fprintf(out, "%s - %s\n", cmd.Name, cmd.Desc)
	fmt.Fprintf(out, "%s\n", helpCommandShort(cmd))
	fmt.Fprintf(out, "\n")

	helpOptionsRecurseRev(cmd)

	if len(cmd.Children) > 0 {
		fmt.Fprintf(out, "Available subcommands:\n")
		for _, v := range cmd.Children {
			fmt.Fprintf(out, "  %-12s %s\n", v.Name, v.Desc)
		}
		fmt.Fprintf(out, "\n")
	}

	if cmd.Handler == nil {
		fmt.Fprintf(out, "For help information run:\n")
		fmt.Fprintf(out, "  '%s help' .. '%s <commands>* help' .. '%s [commands]* help [subcommand]*'\n",
			cmd.GetNameTop(), cmd.GetNameTop(), cmd.GetNameTop())
		fmt.Fprintf(out, "\n")
	}
}

func helpCommandShort(cmd *Command) string {
	var params []string

	for _, option := range cmd.Options {
		params = append([]string{helpCommandShortOption(option)}, params...)
	}

	params = append(params, cmd.Name)

	if cmd.Parent != nil {
		params = append(params, helpCommandShort(cmd.Parent))
	}

	for i, j := 0, len(params)-1; i < j; i, j = i+1, j-1 {
		params[i], params[j] = params[j], params[i]
	}

	return strings.Join(params, " ")
}

func helpCommandShortOption(option *Option) string {
	var optstr string

	if !option.required {
		optstr += "["
	}

	if option.param {
		optstr += "--" + option.name + " <" + option.name + ">"
	} else {
		optstr += "-" + option.name
	}

	if !option.required {
		optstr += "]"
	}

	return optstr
}

func helpOptionsRecurseRev(cmd *Command) {
	if cmd.Parent != nil {
		helpOptionsRecurseRev(cmd.Parent)
	}

	helpOptions(cmd)
}

func helpOptions(cmd *Command) {
	if len(cmd.Options) == 0 {
		return
	}

	fmt.Printf("%s options:\n", cmd.GetNameChain())
	for _, option := range cmd.Options {
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
