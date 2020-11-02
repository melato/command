package main

import (
	"fmt"

	"melato.org/command"
)

type N float32

type SimpleDemo struct {
	S string
	I int
	U uint32
	C complex64
	F N
}

func (t *SimpleDemo) Run(args []string) error {
	fmt.Printf("S=%s I=%d N=%f C=%v\n", t.S, t.I, t.F, t.C)
	return nil
}

func main() {
	t := &SimpleDemo{}
	var cmd command.SimpleCommand
	command.Main(cmd.Flags(t).Method(t.Run))
}
