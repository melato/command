package main

import (
	"errors"
	"fmt"

	"melato.org/command"
)

// Flags of various types
type App struct {
	S       string `name:"s,s-flag" usage:"string flag with two names"`
	B       bool   // bool flag with default name, no usage
	IntFlag int    `name:"i" usage:"int flag"`
	F       `name:"f" usage:"aliased float32"`
	Sub1    Sub  // struct
	Sub2    *Sub `name:"sub2" usage:"sub-2: "` // pointer to struct, with prefix name
	Sub3    Sub  `name:""`                     // no flags
}

type F float32

type Sub struct {
	X string `usage:"X"`
	Y string `usage:"Y"`
}

// Subcommand with additional flag
type Hello struct {
	App    *App `name:""` // flags are specified by the parent command
	Prefix string
}

// Initialize some flags (optional).
func (t *App) Init() error {
	t.S = "s-default"
	t.Sub2.Y = "y-default"
	return nil
}

// Check if the flags defined by the user are ok, before running the command.
func (t *App) Configured() error {
	if t.S == "" {
		return errors.New("missing -s")
	}
	return nil
}

func (t *App) printFlags() error {
	fmt.Printf("s: %s\n", t.S)
	fmt.Printf("b: %v\n", t.B)
	fmt.Printf("i: %d\n", t.IntFlag)
	fmt.Printf("f: %f\n", t.F)
	fmt.Printf("sub.x: %s\n", t.Sub1.X)
	return nil
}

func (t *App) doError() error {
	return errors.New("here's an error")
}

func (t *Hello) Init() error {
	t.Prefix = "hello"
	return nil
}

func (t *Hello) Hello(args []string) error {
	t.App.printFlags()
	for _, arg := range args {
		fmt.Printf("%s %s\n", t.Prefix, arg)
	}
	return nil
}

// main has the only dependency on command.  It can be in a file on its own.
func main() {
	app := &App{Sub2: &Sub{"x2", "y2"}}
	var cmd command.SimpleCommand
	cmd.Flags(app)
	cmd.Command("flags").RunMethodE(app.printFlags).Short("demonstrate flag processing")
	cmd.Command("error").RunMethodE(app.doError).Short("demonstrate error handling")

	hello := &Hello{App: app}
	cmd.Command("hello").Flags(hello).RunMethodArgs(hello.Hello).Short("sub-command with args and additional flags").Use("[arg]...")
	command.Main(&cmd)
}
