package main

import (
	"errors"

	"melato.org/command"
)

type App struct {
	Trace bool `name:"t" usage:"trace"`
}

/* Optional methods:
func (t *App) Init() error {
	return nil
}

func (t *App) Configured() error {
	return nil
}

*/

func (t *App) Run(args []string) error {
	if len(args) != 2 {
		return errors.New("")
	}
	return nil
}

func main() {
	var cmd command.SimpleCommand
	var app App
	cmd.Flags(&app).RunMethodArgs(app.Run).Use("<arg1> <arg2>").Short("example command with two args")
	command.Main(&cmd)
}
