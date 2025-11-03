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

func addSubcommand(cmd *command.SimpleCommand, name string) {
	sub := &Flags{Name: name}
	cmd.Command(name).Flags(sub).RunFunc(sub.Run)
}

func main() {
	cmd := &command.SimpleCommand{}
	cmd.Flags(&Flags{Name: "top"})
	addSubcommand(cmd, "sub1")
	addSubcommand(cmd, "sub2")
	command.Main(cmd)
}
