package command

import (
	"errors"
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

type command interface {
	/** run the command.
	 */
	run(args []string) error

	/** Called before any other method, as a constructor
	It may set default values, which are shown in the usage help.
	*/
	init() error

	/** Called after flags have been applied and there are no flag errors, and the help flag is not present.
	Do not setup sub-commands here, because they may be used from the help, before Configured() is called.*/
	configured() error

	enabledConfig() bool

	cleanup() error

	/** Returns usage information
	 */
	usage() *Usage

	flags() interface{}

	// Return the subcommands.  Flag values have been applied.
	// Configured() may or may not have been called.
	Commands() map[string]*SimpleCommand
}

type commandInfo struct {
	Name           string
	Command        command
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
	err := t.Command.init()
	if err != nil {
		return err
	}
	// extractFlags must be called after Command.Init(),
	// because Command.Init may create flags by assigning values to struct pointers
	t.Flags = extractFlags(t.Command.flags(), &flagPrefix{})
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

func createCommandInfo(name string, cmd command) *commandInfo {
	c := &commandInfo{Name: name, Command: cmd}
	usage := cmd.usage()
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

func showUsage(levels []*commandInfo, commands map[string]*SimpleCommand) {
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

		if len(u.Examples) > 0 {
			fmt.Println()
			fmt.Println("Examples:")
			for _, ex := range u.Examples {
				fmt.Printf("  %s %s\n", levels[0].Name, ex)
			}
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
		var nameLen int
		for _, ci := range ar {
			w := len(ci.Name)
			if w > nameLen {
				nameLen = w
			}
		}
		for _, ci := range ar {
			fmt.Printf("  %-*s  %s\n", nameLen, ci.Name, ci.Usage.Short)
		}
	}
}

func runCommand(name string, cmd command, args []string, ancestors []*commandInfo) error {
	var err error
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
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
	commands := cmd.Commands()

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
				fmt.Fprintf(os.Stderr, "no such command: %s\n", name2)
				showUsage(ancestors, commands)
				os.Exit(1)
			}
		} else {
			showUsage(ancestors, commands)
			os.Exit(0)
		}
	} else {
		// call all command-chain Configured() methods just before Run()
		if cmd.enabledConfig() {
			for i, a := range ancestors {
				err := a.Command.configured()
				if err != nil {
					cleanup(ancestors[0 : i+1])
					return err
				}
			}
		}
		err := cmd.run(args2)
		cleanup(ancestors)
		if err != nil && err.Error() == "" {
			u := cmd.usage()
			if u != nil && u.Use != "" {
				err = errors.New("usage: " + u.Use)
			} else {
				err = errors.New("wrong usage")
			}
		}
		return err
	}
	return nil
}

func Main(cmd command) {
	name := filepath.Base(os.Args[0])
	err := runCommand(name, cmd, os.Args[1:], nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cleanup(commands []*commandInfo) {
	for j := len(commands) - 1; j >= 0; j-- {
		if err2 := commands[j].Command.cleanup(); err2 != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err2)
		}
	}
}
