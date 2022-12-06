package main

import (
	_ "embed"

	"melato.org/command"
	"melato.org/command/examples/cli"
	"melato.org/command/usage"
)

//go:embed usage.yaml
var usageData []byte

func main() {
	var cmd command.SimpleCommand

	formatCmd := cmd.Command("format")

	var formatFloat cli.FormatFloat
	formatCmd.Command("float").Flags(&formatFloat).RunFunc(formatFloat.Format)

	var formatTime cli.FormatTime
	formatCmd.Command("time").Flags(&formatTime).RunFunc(formatTime.Format)

	var add cli.Add
	cmd.Command("add").RunFunc(add.Float)

	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
