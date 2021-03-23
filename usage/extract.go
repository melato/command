package usage

import (
	"fmt"

	"gopkg.in/yaml.v2"
	"melato.org/command"
)

// Extract copies the usage fields from a command.SimpleCommand, recursively.
// It can be used to generate an external usage file from hardcoded usage strings,
// so you can then remove the hardcoded usage and replace it with the yaml file.
func Extract(cmd *command.SimpleCommand) Usage {
	var u Usage
	u.Usage = cmd.Usage
	for name, sub := range cmd.Commands() {
		if u.Commands == nil {
			u.Commands = make(map[string]*Usage)
		}
		su := Extract(sub)
		u.Commands[name] = &su
	}
	return u
}

func ExtractToYaml(cmd *command.SimpleCommand) error {
	data, err := yaml.Marshal(Extract(cmd))
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
