package main

import (
	"fmt"

	"melato.org/command"
)

type F1 struct {
	A string
}

type F2 struct {
	B int
}

func (t *F1) F1() error {
	fmt.Printf("F1 A=%s\n", t.A)
	return nil
}

func (t *F2) F2(args []string) error {
	fmt.Printf("F2 B=%d\n", t.B)
	for _, arg := range(args) {
		fmt.Printf("%s\n", arg)
	}
	return nil
}

type Level2 struct {
	C string
}

func (t *Level2) X() {
	fmt.Printf("x c=%s\n", t.C)
}

func (t *Level2) Y() {
	fmt.Printf("y c=%s\n", t.C)
}

func main() {	
	var cmd command.SimpleCommand
	var f1 F1
	var f2 F2
	var v2 Level2
	cmd.Command("f1").Flags(&f1).RunMethodE(f1.F1).Short("Demo string option")
	cmd.Command("f2").Flags(&f2).RunMethodArgs(f2.F2).Use("[arg]...").Short("Demo options + args")
	level2 := cmd.Command("sub").Flags(&v2)
	level2.Command("x").Flags(&v2).RunMethod(v2.X)
	level2.Command("y").Flags(&v2).RunMethod(v2.Y)
	command.Main(&cmd)
}
