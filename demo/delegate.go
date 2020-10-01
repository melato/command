package main

import (
	"fmt"

	"melato.org/command"
)

type S1 struct {
	B string `usage:"B"`
}

type S2 struct {
	C string
}

func (t *S2) String() string {
	return t.C
}

type Demo struct {
	command.Base
	A string
	X S1 `name:"x" usage:"X"`
	Y S1 `name:"y" usage:"Y"`
	//Z *S2
	Z interface{}
}

func (t *Demo) Init() error {
	fmt.Println("Init()")
	t.Z = &S2{C: "generic"}
	return nil
}
func (t *Demo) Run(args []string) error {
	fmt.Println("Demo", t.A, t.X.B, t.Y.B, t.Z)
	return nil
}

func main() {
	command.Main(&Demo{})
}
