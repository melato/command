package main

import (
	_ "embed"
	"fmt"

	"example.org/example/cli"
	"melato.org/command"
	"melato.org/command/usage"
)

//go:embed usage.yaml
var usageData []byte

//go:embed version
var version string

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
	stringCmd.Command("sprintf").RunFunc(str.Sprintf)

	var re cli.Regexp
	reCmd := cmd.Command("regexp")
	reCmd.Flags(&re)
	reCmd.Command("split").RunFunc(re.Split)
	reCmd.Command("submatch").RunFunc(re.FindStringSubmatch)
	reCmd.Command("find").RunFunc(re.FindAllString)

	cmd.Command("version").NoConfig().RunFunc(func() { fmt.Printf("%s\n", version) })
	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
