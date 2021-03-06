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
	if len(args) == 0 {
		// an empty error is a usage error.
		// Command will replace an empty error with usage error that includes the command's Use message.
		return errors.New("")
	}
	t.App.printFlags()
	for _, arg := range args {
		fmt.Printf("%s %s\n", t.Prefix, arg)
	}
	return nil
}

func (t *App) Two(a, b string) {
	fmt.Println("two:", a, b)
}

func (t *App) OnePlus(a string, b ...string) {
	fmt.Println("one+:", a, b)
}

func (t *App) Mixed(a string, n int) error {
	fmt.Printf("a=%s n=%d\n", a, n)
	return nil
}

// main has the only dependency on command.  It can be in a file on its own.
func main() {
	app := &App{Sub2: &Sub{"x2", "y2"}}
	var cmd command.SimpleCommand
	cmd.Flags(app).Short("demonstrate command usage").
		Long(`This demo shows how to use commands and subcommands that populate struct fields from command line flags
and execute functions with various signatures.`)
	cmd.Command("flags").RunFunc(app.printFlags).Short("demonstrate flag processing")
	cmd.Command("error").RunFunc(app.doError).Short("demonstrate error handling")

	hello := &Hello{App: app}
	cmd.Command("hello").Flags(hello).RunFunc(hello.Hello).Short("sub-command with args and additional flags").Use("arg...")
	cmd.Command("two").RunFunc(app.Two).Short("two strings")
	cmd.Command("one+").RunFunc(app.OnePlus).Short("1+ args").Use("<a> [b]...")
	cmd.Command("mixed").RunFunc(app.Mixed).Use("<string> <int>").Short("args with mixed types")
	command.Main(&cmd)
}
