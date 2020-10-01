package main

import (
	"fmt"

	"melato.org/export/command"
)

type SimpleDemo struct {
	A string
	B int
}

func (t *SimpleDemo) Run(args []string) error {
	fmt.Printf("A=%s B=%d\n", t.A, t.B)
	return nil
}

func main() {
	t := &SimpleDemo{}
	var cmd command.SimpleCommand
	command.Main(cmd.Flags(t).Method(t.Run))
}
