package clicommand

type Handler func(*Data) (err error)

type Data struct {
	Cmd     *Command
	Options map[string]string
	Params  []string
}
