package main

import (
	"fmt"

	"melato.org/command"
)

/** RootCommand - demonstrate Command method calling */
type RootCommand struct {
	command.Base
}

func main() {
	command.Main(&RootCommand{})
}

func (t *RootCommand) Init() error {
	fmt.Println("Root.Init()")
	commands := t.Commands()
	commands["c1"] = &Command1{}
	commands["c2"] = &Command2{}
	return nil
}

func (t *RootCommand) Configured() error {
	fmt.Println("Root.Configured()")
	return nil
}

func (t *RootCommand) Commands() map[string]command.Command {
	fmt.Println("Root.Commands()")
	return t.Base.Commands()
}

func (t *RootCommand) Run(args []string) error {
	fmt.Println("Root.Run() should never be called, when there are sub-commands ")
	return nil
}

type Command1 struct {
	command.Base
}

func (t *Command1) Run(args []string) error {
	fmt.Println("c1 Run()")
	return nil
}

type Command2 struct {
	command.Base
}

func (t *Command2) Run(args []string) error {
	fmt.Println("c2 Run()")
	return nil
}
