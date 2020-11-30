package command

import (
	"errors"
	"strings"
)

type Init interface {
	Init() error
}

type Configured interface {
	Configured() error
}

/** SimpleCommand - a command with no flags that can be defined in a single line with the Add() method. */
type SimpleCommand struct {
	Base
	runMethod    func([]string) error `name:"-"`
	configured   func() error         `name:"-"`
	U            Usage                `name:"-"`
	CommandFlags interface{}
}

func (t *SimpleCommand) Use(s string) *SimpleCommand {
	t.U.Use = s
	return t
}
func (t *SimpleCommand) Short(s string) *SimpleCommand {
	t.U.Short = s
	return t
}
func (t *SimpleCommand) Long(s string) *SimpleCommand {
	t.U.Long = s
	return t
}
func (t *SimpleCommand) Example(s string) *SimpleCommand {
	t.U.Example = s
	return t
}

/** Specify a struct that defines command flags.  The structure is filled with the parsed flags.
If the struct implements the Init() or Configured() interfaces, the Init() and/or Configured()
methods are called as the command Init(), Configured().
*/
func (t *SimpleCommand) Flags(flags interface{}) *SimpleCommand {
	t.CommandFlags = flags
	return t
}

/** Specify which method to run when executing this command */
func (t *SimpleCommand) RunMethodArgs(method func([]string) error) *SimpleCommand {
	t.runMethod = method
	return t
}

/** Specify which method to run when executing this command */
func (t *SimpleCommand) RunMethod(method func()) *SimpleCommand {
	return t.RunMethodArgs(func(args []string) error {
		if len(args) != 0 {
			return errors.New("unrecognized arguments: " + strings.Join(args, " "))
		}
		method()
		return nil
	})
}

/** Specify which method to run when executing this command */
func (t *SimpleCommand) RunMethodE(method func() error) *SimpleCommand {
	return t.RunMethodArgs(func(args []string) error {
		if len(args) != 0 {
			return errors.New("unrecognized arguments: " + strings.Join(args, " "))
		}
		return method()
	})
}

/** Specify which method to run when executing this command */
func (t *SimpleCommand) Method(method func([]string) error) *SimpleCommand {
	return t.RunMethodArgs(method)
}

/** Specify which method to run after configuration.
  Consider using Flags() is called with a Configured() interface instead.
*/
func (t *SimpleCommand) ConfiguredMethod(method func() error) *SimpleCommand {
	t.configured = method
	return t
}

/** Add a subcommand
 */
func (t *SimpleCommand) Command(name string) *SimpleCommand {
	c := &SimpleCommand{}
	t.Commands()[name] = c
	return c
}

/** Add a subcommand
 */
func (t *SimpleCommand) Add(name string, method func([]string) error, short string) *SimpleCommand {
	return t.Command(name).RunMethodArgs(method).Short(short)
}

/** Add a subcommand
 */
func (t *SimpleCommand) AddCommand(name string, cmd Command) {
	t.Commands()[name] = cmd
}

func (t *SimpleCommand) Run(args []string) error {
	return t.runMethod(args)
}

func (t *SimpleCommand) Init() error {
	f, ok := t.CommandFlags.(Init)
	if ok {
		return f.Init()
	}
	return nil
}

func (t *SimpleCommand) Configured() error {
	if t.configured != nil {
		return t.configured()
	}
	f, ok := t.CommandFlags.(Configured)
	if ok {
		return f.Configured()
	}
	return nil
}

func (t *SimpleCommand) Usage() *Usage {
	return &t.U
}
