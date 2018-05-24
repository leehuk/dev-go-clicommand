package clicommand

// An Option represents a defined option parameter that can be specified
// on the command line when the program is run.
//
// Options can be attached to multiple Command objects.
//
// Options either have, or do not have parameters and the parser detects
// the difference by whether they are specified with a double dash prefix
// or a single dash prefix.  Options with a double dash prefix (e.g. --foo)
// will have the next field of the command string taken as its parameter.
//
// The parser will always perform case-insensitive matches on option names,
// but will also match the parameter type.  If a non-parameter (e.g. -bar)
// is specified, but the Option has been added as a parameter type (e.g. --bar)
// the parser will treat it as an unknown option.
type Option struct {
	// name Name of option, stored without its dashes prefix.
	Name string

	// desc Description of option
	Desc string

	// param Controls whether this option takes parameters or not.
	//
	// For simplicity, all options that take parameters have a double dash
	// prefix, whilst options without parameters have a single dash prefix.
	// E.g.
	//   --option1 <required parameter>
	//   -option2 // no paramater
	Param bool

	// required Controls whether this option must be supplied or not.
	// If a parameter is marked as required, the parser will automatically
	// detect it is not supplied and return an error.
	Required bool

	// parents Array of pointers which this Option is bound to
	Parents []*Command
}

// NewOption creates a new Option object with the given name and desc, but does
// not bind it within the tree.  This is generally useful when creating a
// generic Option, which needs to be bound to multiple Command objects.
//
// The param field is used to determine whether this option takes an
// additional parameter after it.
func NewOption(name string, desc string, param bool) *Option {
	opt := &Option{
		Name:  name,
		Desc:  desc,
		Param: param,
	}

	return opt
}

// BindCommand binds an Option to the given Command object, so it is
// available to be specified for that command, and all child commands.
func (o *Option) BindCommand(cmd *Command) {
	o.Parents = append(o.Parents, cmd)
	cmd.Options = append(cmd.Options, o)
}

// UnbindCommand unbinds an Option from the given Command object, at which
// point it is no longer available for that command or its children.
func (o *Option) UnbindCommand(cmd *Command) {
	var newparents []*Command
	var newoptions []*Option

	for _, cmdi := range o.Parents {
		if cmdi != cmd {
			newparents = append(newparents, cmdi)
		}
	}

	for _, opt := range cmd.Options {
		if opt != o {
			newoptions = append(newoptions, opt)
		}
	}

	o.Parents = newparents
	cmd.Options = newoptions
}

// GetRequired returns whether this Option must be specified.  This requirement
// only applies to Options that are directly on the path between the edge Command
// and the root.
func (o *Option) GetRequired() bool {
	return o.Required
}

// SetRequired marks the Option so it must be specified.  This requirement only
// applies to Options that are directly on the path between the edge Command
// and the root.
func (o *Option) SetRequired() *Option {
	o.Required = true
	return o
}

// GetParents returns the parents Command objects of an Option
func (o *Option) GetParents() []*Command {
	return o.Parents
}
