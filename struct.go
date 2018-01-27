package clicommand

type Data struct {
	Cmd     *Command
	Options map[string]string
	Params  []string
}
