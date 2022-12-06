package command

import (
	"errors"
	"strings"
)

// Init is an optional interface for flags objects.
// See SimpleCommand.Flags
type Init interface {
	// Init  is called before any other method, as a constructor.
	//
	// It typically sets default values for flags, which are also shown in the usage help.
	//
	// If it returns an error, the command is not run, and the error is reported as the program's error.
	Init() error
}

// Configured is an optional interface for flags objects that have a Configured() method.
// See SimpleCommand.Flags
type Configured interface {
	// Configured is called just before running the command, after flags have been set.
	//
	// The Configured() methods of any ancestor commands, are called in order, from the root command to the command that is about to run.
	//
	// If it returns an error, the command does not run, and the error is reported as the program's error.
	Configured() error
}

// Closer is an optional interface for flags objects that have a Close() method.
// See SimpleCommand.Flags
// Closer is the same interface as io.Closer().
type Closer interface {
	// Close is called after the command runs, to do any cleanup.
	//
	// Close will be called even if Configured() or the command return an error.
	//
	// It is meant to be used in situations where a Configured() method creates a temporary file, which should be deleted when the program exits.
	Close() error
}

// SimpleCommand defines a CLI command.
//
// Command flags are specified by Flags().
//
// The method to run is specified by one of the RunMethod... calls.
//
// Most methods return the command, so they can be chained together to configure the command.
type SimpleCommand struct {
	subcommands  map[string]*SimpleCommand
	runMethod    func([]string) error
	Usage        Usage
	commandFlags interface{} // The argument that was passed to the Flags() method.  This is meant for internal use.
	noConfig     bool
}

// A generic representation of the command-line arguments, without any options, e.g. "<arg1> <arg2>"
func (t *SimpleCommand) Use(commandLineUsage string) *SimpleCommand {
	t.Usage.Use = commandLineUsage
	return t
}

// A one-line description of the command, shown in lists of commands
func (t *SimpleCommand) Short(shortDescription string) *SimpleCommand {
	t.Usage.Short = shortDescription
	return t
}

// A longer description shown in the help for a single command
func (t *SimpleCommand) Long(longDescription string) *SimpleCommand {
	t.Usage.Long = longDescription
	return t
}

// A usage example, shown in the help for a single command.  May be called multiple times to add examples.
func (t *SimpleCommand) Example(example string) *SimpleCommand {
	t.Usage.Examples = append(t.Usage.Examples, example)
	return t
}

// Flags specifies a pointer to a struct that defines command flags.
//
// The struct fields are set with the parsed flags.
//
// If flags implements the Init, Configured, or Closer interfaces,
// flags.Init(), flags.Configured(), or flags.Close() are called as specified in the interface documentation.
func (t *SimpleCommand) Flags(flags interface{}) *SimpleCommand {
	t.commandFlags = flags
	return t
}

func (t *SimpleCommand) flags() interface{} {
	return t.commandFlags
}

// Specify the method to run when executing this command.  The command arguments are passed to the method.
func (t *SimpleCommand) RunMethodArgs(method func([]string) error) *SimpleCommand {
	t.runMethod = method
	return t
}

// Specify the method to run when executing this command.  It is an error if any arguments are passed to the command.
func (t *SimpleCommand) RunMethodE(method func() error) *SimpleCommand {
	return t.RunMethodArgs(func(args []string) error {
		if len(args) != 0 {
			return errors.New("unrecognized arguments: " + strings.Join(args, " "))
		}
		return method()
	})
}

// RunMethod is like RunMethodE, but does not return an error.
func (t *SimpleCommand) RunMethod(method func()) *SimpleCommand {
	return t.RunMethodE(func() error {
		method()
		return nil
	})
}

// RunFunc specifies the function to run when executing this command.
// It is like the RunMethod* methods, but it uses reflection to match the function arguments to the provided arguments.
// The command arguments are passed to the function.
// fn must return either 0 values or one value that is assignable to error
// It may have any number of arguments of any primitive type (that can be parsed from a string)
func (t *SimpleCommand) RunFunc(fn interface{}) *SimpleCommand {
	return t.RunMethodArgs(wrapFunc(fn))
}

// Commands returns the list of subcommands
// This is not normally needed for adding sub-commands, because sub-commands can be added via Command().
// Commands is exported for traversing the tree of commands, for documentation, or testing.
func (t *SimpleCommand) Commands() map[string]*SimpleCommand {
	if t.subcommands == nil {
		t.subcommands = make(map[string]*SimpleCommand)
	}
	return t.subcommands
}

// AddCommand adds another command as a subcommand.  It returns the sub-command.
func (t *SimpleCommand) AddCommand(name string, c *SimpleCommand) {
	t.Commands()[name] = c
}

// Command creates a subcommand and adds it to this command.  It returns the sub-command.
func (t *SimpleCommand) Command(name string) *SimpleCommand {
	c := &SimpleCommand{}
	t.Commands()[name] = c
	return c
}

func (t *SimpleCommand) run(args []string) error {
	if t.runMethod != nil {
		return t.runMethod(args)
	}
	// there is no run method, so we do nothing.
	// this is used in our demo programs, so we don't want to crash
	return nil
}

func (t *SimpleCommand) init() error {
	f, ok := t.commandFlags.(Init)
	if ok {
		return f.Init()
	}
	return nil
}

// Disable calling of Configured() for flags of this or any ancestor command.
//
// Use for special commands like "version" that should not require the user to enter correct options. */
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

func (t *SimpleCommand) cleanup() error {
	f, ok := t.commandFlags.(Closer)
	if ok {
		return f.Close()
	}
	return nil
}

func (t *SimpleCommand) usage() *Usage {
	return &t.Usage
}
