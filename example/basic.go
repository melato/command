package main

import (
	_ "embed"

	"fmt"

	"melato.org/command"
	"melato.org/command/usage"
)

type App struct {
	S string `name:"s" usage:"example flag"`
}

func (t *App) Init() error {
	t.S = "hello"
	return nil
}

func (t *App) Run() {
	fmt.Printf("s=%s\n", t.S)
}

//go:embed usage.yaml
var usageData []byte

func main() {
	var cmd command.SimpleCommand

	var app App
	cmd.Command("run").Flags(&app).RunFunc(app.Run)

	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
