package main

import (
	"fmt"

	"melato.org/command"
)

type Hello struct {
	Prefix string
}

// Init is optional.  Used to initialize flags
func (t *Hello) Init() error {
	t.Prefix = "hello "
	return nil
}

func (t *Hello) Hello(args []string) {
	for _, name := range args {
		fmt.Println(t.Prefix + name)
	}
}

func main() {
	cmd := &command.SimpleCommand{}
	hello := &Hello{}
	cmd.Flags(hello).RunFunc(hello.Hello).Use("<name>...").Short("add a greeting to a name")
	command.Main(cmd)
}
