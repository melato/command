package main

import (
	_ "embed"
	"fmt"

	"melato.org/command"
	"melato.org/command/example/cli"
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

	cmd.Command("version").NoConfig().RunFunc(func() { fmt.Printf("%s\n", version) })
	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
