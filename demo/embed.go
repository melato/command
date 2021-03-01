package main

import (
	_ "embed"

	"melato.org/command"
	"melato.org/command/usage"
)

//go:embed embed.yaml
var usageData []byte

// demonstrates how to set command usage from embedded file
func main() {
	cmd := &command.SimpleCommand{}
	// we don't need Run functions here, so we don't use them
	cmd.Command("a")
	b := cmd.Command("b")
	b.Command("b1")
	b.Command("b2")
	b.Command("b3")
	cmd.Command("c")

	_ = usage.ApplyEnv(cmd, "USAGE_FILE") || usage.ApplyYaml(cmd, usageData)

	command.Main(cmd)
}
