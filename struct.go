package clicommand

type Command struct {
	name      string
	desc      string
	handler   Handler
	parent    *Command
	children  []*Command
	args      []*Arg
	callbacks []Handler
}

type Arg struct {
	name     string
	desc     string
	descx    []string
	param    bool
	required bool
}

type Handler func(*Data) (err error)

type Data struct {
	Cmd     *Command
	Options map[string]string
	Params  []string
}
