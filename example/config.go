package main

import (
	"fmt"

	"melato.org/command"
)

type Flags struct {
	Name string `name:"-"`
}

func (t *Flags) Run() {
	fmt.Printf("%s.Run()\n", t.Name)
}

func (t *Flags) Init() error {
	fmt.Printf("%s.Init()\n", t.Name)
	return nil
}

func (t *Flags) Configured() error {
	fmt.Printf("%s.Configured()\n", t.Name)
	return nil
}

func main() {
	var cmd command.SimpleCommand
	cmd.Flags(&Flags{Name: "top"})
	sub := &Flags{Name: "sub"}
	cmd.Command("sub").Flags(sub).RunFunc(sub.Run)
	command.Main(&cmd)
}
