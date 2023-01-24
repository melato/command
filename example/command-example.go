package main

import (
	_ "embed"
	"fmt"

	"example.org/command/cli"
	"melato.org/command"
	"melato.org/command/usage"
)

//go:embed command-example.yaml
var usageData []byte

//go:embed cli/flags.yaml
var flagsUsage []byte

//go:embed cli/funcs.yaml
var funcsUsage []byte

//go:embed version
var version string

func FlagsCommand() *command.SimpleCommand {
	app := &cli.Flags{Sub2: &cli.EmbeddedType{"x2", "y2"}}
	var cmd command.SimpleCommand
	cmd.Flags(app)
	cmd.Command("print").RunFunc(app.PrintFlags)

	sub := &cli.AdditionalFlags{Types: app}
	cmd.Command("sub").Flags(sub).RunFunc(sub.Run)
	return &cmd
}

func FuncsCommand() *command.SimpleCommand {
	var cmd command.SimpleCommand
	var funcs cli.Funcs
	cmd.Command("error").RunFunc(funcs.DoError)
	cmd.Command("strings").RunFunc(funcs.Strings)
	cmd.Command("floats").RunFunc(funcs.Floats)
	cmd.Command("mixed").RunFunc(funcs.Mixed)
	usage.ApplyEnv(&cmd, "USAGE", funcsUsage)
	return &cmd
}

func main() {
	var cmd command.SimpleCommand

	formatCmd := cmd.Command("format")

	var formatFloat cli.FormatFloat
	formatCmd.Command("float").Flags(&formatFloat).RunFunc(formatFloat.Format)

	var formatTime cli.FormatTime
	formatCmd.Command("time").Flags(&formatTime).RunFunc(formatTime.Format)

	var add cli.Add
	cmd.Command("add").RunFunc(add.Float)

	var str cli.Strings
	stringCmd := cmd.Command("string")
	stringCmd.Command("join").Flags(&str).RunFunc(str.Join)
	stringCmd.Command("split").Flags(&str).RunFunc(str.Split)

	var sprintf cli.Sprintf
	cmd.Command("sprintf").Flags(&sprintf).RunFunc(sprintf.Sprintf)

	var re cli.Regexp
	reCmd := cmd.Command("regexp")
	reCmd.Flags(&re)
	reCmd.Command("split").RunFunc(re.Split)
	reCmd.Command("submatch").RunFunc(re.FindStringSubmatch)
	reCmd.Command("find").RunFunc(re.FindAllString)

	cmd.AddCommand("flags", FlagsCommand())
	cmd.AddCommand("funcs", FuncsCommand())
	cmd.Command("version").NoConfig().RunFunc(func() { fmt.Printf("%s\n", version) })
	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
