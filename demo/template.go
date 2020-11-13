package main

import (
	"fmt"

	"melato.org/command"
)

type App struct {
}

func (t *App) Run(args []string) error {
	fmt.Println(args)
	return nil
}

func main() {
	var cmd command.SimpleCommand
	var app App
	cmd.Flags(&app).RunMethodArgs(app.Run)
	command.Main(&cmd)
}
