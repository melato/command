/**
  Provides a simple flag-based command line interface (CLI) that uses reflection to define flags and their default values from struct fields.
*/
package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

/** A Command is a struct type, whose fields are used to specify the CLI flags.
Flags are fields that are primitive types (string, int, bool, etc.) or slices of primitive types.
The name of the flag is specified by the "name" tag, or by the name of the field, in lowercase.
The usage string of the flag is specified by the "usage" tag.

The usage string of the whole command is taken from the "usage" tag of an exported command field of type Usage.
*/

type Command interface {
	/** run the command.
	 */
	Run(args []string) error

	/** Called before any other method, as a constructor
	It may set default values, which are shown in the usage help.
	*/
	Init() error

	/** Called after flags have been applied and there are no flag errors, and the help flag is not present.
	Do not setup sub-commands here, because they may be used from the help, before Configured() is called.*/
	Configured() error

	/** Returns usage information
	 */
	Usage() *Usage

	/** Return the subcommands.  Flag values have been applied.
	Configured() may or may not have been called. */
	Commands() map[string]Command
}

/**
Documentation for the command.
*/
type Usage struct {
	/** A generic representation of the command-line arguments, without any options, e.g. "<arg1> <arg2>"  */
	Use string

	/** A one-line description of the command, shown in lists of commands */
	Short string

	/** A longer description shown in the help for a single command */
	Long string

	Example string
}

type commandInfo struct {
	Name           string
	Command        Command
	Usage          Usage
	Flags          []*commandFlag
	extractedFlags bool
	FlagSet        *flag.FlagSet
}

func (t *commandInfo) Init() error {
	if t.extractedFlags {
		return nil
	}
	t.extractedFlags = true
	err := t.Command.Init()
	if err != nil {
		return err
	}
	// extractFlags must be called after Command.Init(),
	// because Command.Init may create flags by assigning values to struct pointers
	t.Flags = extractFlags(t.Command, &flagPrefix{})
	return nil
}

/** should be called after Init() */
func (t *commandInfo) hasOptions() bool {
	return len(t.Flags) > 0
}

type commandInfoSorter []*commandInfo

func (t commandInfoSorter) Len() int {
	return len(t)
}

func (t commandInfoSorter) Swap(i, j int) {
	x := t[i]
	t[i] = t[j]
	t[j] = x
}

func (t commandInfoSorter) Less(i, j int) bool {
	x := t[i]
	y := t[j]
	return strings.Compare(x.Name, y.Name) < 0
}

func createCommandInfo(name string, cmd Command) *commandInfo {
	c := &commandInfo{Name: name, Command: cmd}
	usage := cmd.Usage()
	if usage != nil {
		c.Usage = *usage
	}
	return c
}

/** Create flag.FlagSet flags */
func (c *commandInfo) setFlags(fs *flag.FlagSet) error {
	err := c.Init()
	if err != nil {
		return err
	}
	for _, cf := range c.Flags {
		k := cf.PrimaryNameIndex()
		for i, name := range cf.Names {
			usage := cf.Usage
			if i != k {
				usage = "same as --" + cf.Names[k]
			} else {
				usage = cf.Prefix.ComposeUsage(usage)
			}
			fs.Var(cf.Value, cf.Prefix.ComposeName(name), usage)
		}
	}
	return nil
}

func (u *commandInfo) optionsString(i, n int) string {
	var options string
	if i == 0 {
		options = "Global Options"
	} else if i == n-1 {
		options = "Options"
	} else {
		options = u.Name + " Options"
	}
	return options
}

func showUsage(levels []*commandInfo, commands map[string]Command) {
	var last *commandInfo
	if len(levels) > 0 {
		last = levels[len(levels)-1]
	}
	if last != nil {
		u := last.Usage
		if u.Short != "" {
			fmt.Println(u.Short)
		}
		if u.Long != "" {
			fmt.Println()
			fmt.Println("Description:")
			fmt.Println(u.Long)
		}
		fmt.Println()
		fmt.Println("Usage:")
		var cargs []interface{}
		for _, ci := range levels {
			cargs = append(cargs, ci.Name)
			if ci.hasOptions() {
				cargs = append(cargs, "[options]")
			}
		}
		if len(commands) > 0 {
			cargs = append(cargs, "<command>")
		}
		argsUsage := u.Use
		if argsUsage != "" {
			cargs = append(cargs, argsUsage)
		}

		fmt.Println(cargs...)

		if len(u.Example) > 0 {
			fmt.Println()
			fmt.Println("Examples:")
			fmt.Println(strings.TrimSpace(u.Example))
		}
	}
	for n, i := len(levels), len(levels)-1; i >= 0; i-- {
		u := levels[i]
		if u.hasOptions() {
			fmt.Println()
			fmt.Println(u.optionsString(i, n) + ":")
			u.FlagSet.PrintDefaults()
		}
	}

	if len(commands) > 0 {
		fmt.Println()
		fmt.Println("Available Commands:")
		var ar []*commandInfo
		for name, cmd := range commands {
			ci := createCommandInfo(name, cmd)
			ar = append(ar, ci)
		}
		sort.Sort(commandInfoSorter(ar))
		rows := make([][]string, 0, len(ar))
		for _, ci := range ar {
			rows = append(rows, []string{ci.Name, ci.Usage.Short})
		}
		table := &FixedColumnsTable{Prefix: "  ", Separator: " "}
		table.Print(rows)
	}
}

func runCommand(name string, cmd Command, args []string, ancestors []*commandInfo) error {
	var err error
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	ci := createCommandInfo(name, cmd)
	err = ci.setFlags(fs)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	ci.FlagSet = fs
	// add a help flag
	var help bool
	if fs.Lookup("h") == nil {
		fs.BoolVar(&help, "h", false, "help")
	}
	// parse and apply the flags
	err = fs.Parse(args)

	ancestors = append(ancestors, ci)
	var commands map[string]Command = cmd.Commands()

	if err != nil && err != flag.ErrHelp {
		fmt.Println(err)
		os.Exit(1)
	}

	if help {
		showUsage(ancestors, commands)
		os.Exit(0)
	}

	args2 := fs.Args()
	if len(commands) > 0 {
		if len(args2) > 0 {
			name2 := args2[0]
			cmd2, found := commands[name2]
			if found {
				return runCommand(name2, cmd2, args2[1:], ancestors)
			} else {
				fmt.Println("no such command:", name2)
				showUsage(ancestors, commands)
				os.Exit(1)
			}
		} else {
			showUsage(ancestors, commands)
			os.Exit(0)
		}
	} else {
		// call all command-chain Configured() methods just before Run()
		for _, a := range ancestors {
			err := a.Command.Configured()
			if err != nil {
				return err
			}
		}
		return cmd.Run(args2)
	}
	return nil
}

func ProcessError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Main(cmd Command) {
	name := filepath.Base(os.Args[0])
	ProcessError(runCommand(name, cmd, os.Args[1:], nil))
}
