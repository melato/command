// Package command imlements a command line interface that uses reflection
// to define command flags (options) from the fields of any user-specified struct.
//
// A command has a hierarchy of sub-commands.  Each sub-command can have additional flags.
//
// At each level, optional Init() and Configured() methods can do initialization and validation.
//
// command uses the Go flags package.
package command

import (
	"errors"
	"strings"
)

// Init  is called before any other method, as a constructor.
//
// It typically sets default values for flags, which are also shown in the usage help.
//
// If it returns an error, the command is not run, and the error is reported as the program's error.
type Init interface {
	Init() error
}

// Configured is called just before running the command, after flags have been set.
//
// The Configured() methods of any ancestor commands, are called in order, from the root command to the command that is about to run.
//
// If Configured() returns an error, the command does not run, and the error is reported as the program's error.
type Configured interface {
	Configured() error
}

// SimpleCommand defines a CLI command.
//
// Command flags are specified by Flags().
//
// The method to run is specified by one of the RunMethod... calls.
//
// Most methods return the command, so they can be chained together to configure the command.
type SimpleCommand struct {
	subcommands  map[string]command
	runMethod    func([]string) error
	u            usage
	commandFlags interface{} // The argument that was passed to the Flags() method.  This is meant for internal use.
	noConfig     bool
}

// A generic representation of the command-line arguments, without any options, e.g. "<arg1> <arg2>"
func (t *SimpleCommand) Use(commandLineUsage string) *SimpleCommand {
	t.u.Use = commandLineUsage
	return t
}

// A one-line description of the command, shown in lists of commands
func (t *SimpleCommand) Short(shortDescription string) *SimpleCommand {
	t.u.Short = shortDescription
	return t
}

// A longer description shown in the help for a single command
func (t *SimpleCommand) Long(longDescription string) *SimpleCommand {
	t.u.Long = longDescription
	return t
}

// A usage example, shown in the help for a single command.  May be called multiple times to add examples.
func (t *SimpleCommand) Example(example string) *SimpleCommand {
	t.u.Examples = append(t.u.Examples, example)
	return t
}

// Flags specifies a pointer to a struct that defines command flags.
//
// The struct fields are set with the parsed flags.
//
// If the struct implements the Init() or Configured() interfaces,
// flags.Init() and flags.Configured() are called as specified in the interface documentation.
func (t *SimpleCommand) Flags(flags interface{}) *SimpleCommand {
	t.commandFlags = flags
	return t
}

func (t *SimpleCommand) flags() interface{} {
	return t.commandFlags
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

func (t *SimpleCommand) commands() map[string]command {
	if t.subcommands == nil {
		t.subcommands = make(map[string]command)
	}
	return t.subcommands
}

/** Create a subcommand and add it to the command.  Return the sub-command.
 */
func (t *SimpleCommand) Command(name string) *SimpleCommand {
	c := &SimpleCommand{}
	t.commands()[name] = c
	return c
}

func (t *SimpleCommand) run(args []string) error {
	return t.runMethod(args)
}

func (t *SimpleCommand) init() error {
	f, ok := t.commandFlags.(Init)
	if ok {
		return f.Init()
	}
	return nil
}

/** Disable calling of Configured() for flags of this or any ancestor command.
  Use for special commands like "version" that should not require the user to enter correct options. */
func (t *SimpleCommand) NoConfig() *SimpleCommand {
	t.noConfig = true
	return t
}

func (t *SimpleCommand) enabledConfig() bool {
	return !t.noConfig
}

func (t *SimpleCommand) configured() error {
	f, ok := t.commandFlags.(Configured)
	if ok {
		return f.Configured()
	}
	return nil
}

func (t *SimpleCommand) usage() *usage {
	return &t.u
}
