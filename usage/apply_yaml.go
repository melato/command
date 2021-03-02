// Package usage provides basic utilities for setting command usage from external or embedded files
// see demo/embed.go for an example
package usage

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ApplyYaml Extract usage from Yaml data and applies it to the command hierarchy.
// It prints any errors to stderr.
func ApplyYaml(apply func(Usage), yamlUsage []byte) bool {
	var use Usage
	err := yaml.Unmarshal(yamlUsage, &use)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	apply(use)
	return true
}
