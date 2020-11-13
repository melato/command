package main

import (
	"fmt"

	"melato.org/command"
)

type App struct {
	S string `name:"s" usage:"example flag"`
}

func (t *App) Init() error {
	t.S = "Example"
	return nil
}

func (t *App) Configured() error {
	return nil
}

func (t *App) Run(args []string) error {
	fmt.Println(t.S, args)
	return nil
}

func main() {
	var cmd command.SimpleCommand
	var app App
	cmd.Flags(&app).RunMethodArgs(app.Run)
	command.Main(&cmd)
}
